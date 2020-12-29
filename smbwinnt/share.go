// +build windows

package smbwinnt

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
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
	proc, err := c.dll.FindProc("smb2_get_error")
	if err != nil {
		return fmt.Errorf("failed to resovle error: %v", err)
	}

	pchar, _, _ := proc.Call(c.ptr)
	if pchar != 0 {
		// TODO: Convert uintptr to string?
		return errors.New("unknown smb error")
	}

	return nil
}

// free current smb2_context. Keep it last.
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

// setAuth setup auth parameters.
func (c *context) setAuth(auth *model.Auth) error {
	setAuthProc, err := c.dll.FindProc("smb2_set_authentication")
	if err != nil {
		return err
	}
	setSecModeProc, err := c.dll.FindProc("smb2_set_security_mode")
	if err != nil {
		return err
	}
	setDomainProc, err := c.dll.FindProc("smb2_set_domain")
	if err != nil {
		return err
	}
	setPasswdProc, err := c.dll.FindProc("smb2_set_password")
	if err != nil {
		return err
	}

	switch auth.Type {
	case model.AuthTypeKrb5:
		{
			_, _, err = setAuthProc.Call(c.ptr, 2)
			if err != syscall.Errno(0) {
				return err
			}
			_, _, err = setSecModeProc.Call(c.ptr, 2)
			if err != syscall.Errno(0) {
				return err
			}
		}
	case model.AuthTypeNTLM:
		{
			_, _, err = setAuthProc.Call(c.ptr, 1)
			if err != syscall.Errno(0) {
				return err
			}
			_, _, err = setSecModeProc.Call(c.ptr, 2)
			if err != syscall.Errno(0) {
				return err
			}
		}
	default:
		_, _, err = setSecModeProc.Call(c.ptr, 1)
		if err != syscall.Errno(0) {
			return err
		}
	}

	if len(auth.Domain) > 0 {
		domainPtr, err := syscall.BytePtrFromString(auth.Domain)
		if err != nil {
			return err
		}
		_, _, err = setDomainProc.Call(
			c.ptr,
			uintptr(unsafe.Pointer(domainPtr)),
		)
		if err != syscall.Errno(0) {
			return err
		}
	}

	if len(auth.Password) > 0 {
		passwdPtr, err := syscall.BytePtrFromString(auth.Password)
		if err != nil {
			return err
		}
		_, _, err = setPasswdProc.Call(
			c.ptr,
			uintptr(unsafe.Pointer(passwdPtr)),
		)
		if err != syscall.Errno(0) {
			return err
		}
	}

	return nil
}

// Dial create new session with share. URL must be like "//127.0.0.1/share"
func Dial(urlstr string, auth *model.Auth) (model.Share, error) {
	ctx, err := mkContext()
	if err != nil {
		return nil, fmt.Errorf("failed load dial: %v", err)
	}

	uri, err := url.Parse(urlstr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect share: %v", err)
	}

	// extract share name from uri ()
	parts := strings.Split(uri.Path, "/")
	if len(parts) == 0 {
		return nil, fmt.Errorf("failed to connect share: %v", model.ErrInvalidResourcePath)
	}
	if len(parts[0]) == 0 {
		parts = parts[1:]
	}

	err = ctx.setAuth(auth)
	if err != nil {
		return nil, fmt.Errorf("failed to connect share: %v", err)
	}

	err = conn(ctx, uri.Host, parts[0], auth.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to connect share: %v", err)
	}

	return &share{
		ctx: ctx,
	}, nil
}

// conn creates connection with target share.
func conn(ctx *context, host, share, user string) error {
	if !ctx.ok() {
		return errors.New("failed to dial: context is nilptr")
	}
	proc, err := ctx.dll.FindProc("smb2_connect_share")
	if err != nil {
		return err
	}
	hostPtr, err := syscall.BytePtrFromString(host)
	if err != nil {
		return err
	}
	sharePtr, err := syscall.BytePtrFromString(share)
	if err != nil {
		return err
	}
	userPtr, err := syscall.BytePtrFromString(user)
	if err != nil {
		return err
	}

	_, _, err = proc.Call(
		ctx.ptr,
		uintptr(unsafe.Pointer(hostPtr)),
		uintptr(unsafe.Pointer(sharePtr)),
		uintptr(unsafe.Pointer(userPtr)),
	)
	if err != syscall.Errno(0) {
		return ctx.lastError()
	}

	return nil
}

type share struct {
	ctx *context
}

// ok check ptr state for instance.
func (s *share) ok() bool {
	return s != nil && s.ctx != nil
}

// Echo sends echo request. Impl Share interface.
func (s *share) Echo() error {
	if !s.ok() {
		return errors.New("failed to send echo: context is nilptr")
	}

	proc, err := s.ctx.dll.FindProc("smb2_echo")
	if err != nil {
		return errors.New("failed to send echo: dll func not found")
	}

	n, _, err := proc.Call(s.ctx.ptr)
	if n != 0 || err != syscall.Errno(0) {
		return s.ctx.lastError()
	}

	return nil
}

// OpenFile impl Share interface method.
func (s *share) OpenFile(filename string, mode int) (model.File, error) {
	if !s.ok() {
		return nil, fmt.Errorf("failed to open file: %v", model.ErrContextIsNil)
	}

	return mkFileHandler(s.ctx, filename, mode)
}

// OpenDir impl Share interface method.
func (s *share) OpenDir(path string) (model.Dir, error) {
	if !s.ok() {
		return nil, fmt.Errorf("failed to open dir: %v", model.ErrContextIsNil)
	}

	return mkDirHandler(s.ctx, path)
}

// Close impl Share interface method,
// free all resources for connection.
func (s *share) Close() error {
	if !s.ok() {
		return nil
	}

	proc, err := s.ctx.dll.FindProc("smb2_disconnect_share")
	if err != nil {
		return err
	}
	_, _, err = proc.Call(s.ctx.ptr)
	if err != syscall.Errno(0) {
		return err
	}

	return s.ctx.free()
}
