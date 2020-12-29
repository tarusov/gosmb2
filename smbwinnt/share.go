// +build windows

package smbwinnt

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"

	"github.com/tarusov/gosmb2/model"
)

// context represent app context and dll handlers.
type context struct {
	dll *syscall.DLL
	ptr uintptr
}

// mkContext creates new smb2_context instance.
func mkContext() (*context, error) {
	dll, err := syscall.LoadDLL("smb2.dll")
	if err != nil {
		return nil, fmt.Errorf("failed load smb2.dll: %v", err)
	}
	proc, err := dll.FindProc("smb2_init_context")
	if err != nil {
		return nil, fmt.Errorf("failed to find smb2_init_context: %v", err)
	}
	ptr, _, err := proc.Call()
	if ptr == 0 || err != syscall.Errno(0) {
		return nil, err
	}

	return &context{
		dll: dll,
		ptr: ptr,
	}, nil
}

func (c *context) ok() bool {
	return c.dll != nil && c.ptr != 0
}

// lastError returns last smb2 error from smb2_context in go format.
func (c *context) lastError() error {
	if !c.ok() {
		return errors.New("context is nil")
	}

	proc, err := c.dll.FindProc("smb2_get_error")
	if err != nil {
		return err
	}

	n, _, _ := proc.Call(c.ptr)

	var strval string
	ptr := unsafe.Pointer(n)
	if ptr != nil {
		fmt.Println((*string)(ptr))
	}

	return errors.New(strval)
}

func (c *context) setUsername(s string) {
	//proc, err :=
}

// Free current smb2_context.
func (c *context) free() error {
	if !c.ok() {
		return nil
	}
	proc, err := c.dll.FindProc("smb2_destroy_context")
	if err != nil {
		return fmt.Errorf("failed to find smb2_destroy_context: %v", err)
	}

	_, _, err = proc.Call(c.ptr)
	if err != syscall.Errno(0) {
		return err
	}

	return c.dll.Release()
}

// Dial create new session with share. URL must be like "//127.0.0.1/share"
func Dial(urlstr string, auth *model.Auth) (model.Share, error) {
	ctx, err := mkContext()
	if err != nil {
		return nil, fmt.Errorf("failed load dial: %v", err)
	}

	return &share{
		ctx: ctx,
	}, nil
}

type share struct {
	ctx *context
	ptr uintptr
}

// Echo sends echo request. Impl Share interface.
func (s *share) Echo() error {
	if !s.ctx.ok() {
		return errors.New("failed to send echo: context is nilptr")
	}

	proc, err := s.ctx.dll.FindProc("smb2_echo")
	if err != nil {
		return errors.New("failed to send echo: dll func not found")
	}

	n, _, err := proc.Call(s.ctx.ptr)
	if n != 0 || err != syscall.Errno(0) {
		fmt.Println(n, err)
		return err
	}

	return nil
}

func (s *share) OpenFile(filename string, mode int) (model.File, error) {
	return nil, errors.New("not implemented yet")
}

func (s *share) OpenDir(path string) (model.Dir, error) {
	return nil, errors.New("not implemented yet")
}

// Close impl Share interface method,
// free all resources for connection.
func (s *share) Close() error {
	if s.ptr == 0 || !s.ctx.ok() {
		return nil
	}
	proc, err := s.ctx.dll.FindProc("smb2_disconnect_share")
	if err != nil {
		return err
	}
	_, _, err = proc.Call(s.ctx.ptr, s.ptr)
	if err != syscall.Errno(0) {
		return err
	}

	return s.ctx.free()
}
