// +build windows

package smbwinnt

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/tarusov/gosmb2/model"
)

// file implements model.File.
type file struct {
	ctx  *context
	path string
	ptr  uintptr
	pos  uintptr
}

// mkFileHandler creates new smb2fh instance.
func mkFileHandler(ctx *context, path string, mode int) (*file, error) {
	if !ctx.ok() {
		return nil, fmt.Errorf("failed to open file: %v", model.ErrContextIsNil)
	}

	proc, err := ctx.dll.FindProc("smb2_open")
	if err != nil {
		return nil, err
	}

	pathPtr, err := syscall.BytePtrFromString(path)
	if err != nil {
		return nil, err
	}

	ptr, _, err := proc.Call(
		ctx.ptr,
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(mode),
	)
	if ptr == 0 || err != syscall.Errno(0) {
		return nil, ctx.lastError()
	}

	return &file{
		path: path,
		ptr:  ptr,
		ctx:  ctx,
	}, nil
}

// ok check ptr state for instance.
func (f *file) ok() bool {
	return f != nil && f.ptr != 0 && f.ctx != nil && f.ctx.ptr != 0
}

// close func free current smb2fh data.
func (f *file) close() error {
	if f.ok() {
		proc, err := f.ctx.dll.FindProc("smb2_close")
		if err != nil {
			return err
		}

		_, _, err = proc.Call(f.ctx.ptr, f.ptr)
		if err != syscall.Errno(0) {
			return err
		}
	}

	return nil
}

// Close impl File interface method.
func (f *file) Close() error {
	return f.close()
}

// Stat impl File interface method.
func (f *file) Stat() (os.FileInfo, error) {
	if !f.ok() {
		return nil, fmt.Errorf("failed to get file stat: %v", model.ErrContextIsNil)
	}

	return stat(f.ctx, f.path)
}
