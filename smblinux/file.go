// +build linux

package smblinux

/*
#cgo CFLAGS:  -I./include
#cgo amd64   LDFLAGS: -L./lib/amd64 -lsmb2 -lgssapi_krb5 -lkrb5 -lk5crypto -lkrb5support -lcom_err -ldl -lresolv -lpthread
#cgo 386     LDFLAGS: -L./lib/386   -lsmb2 -lgssapi_krb5 -lkrb5 -lk5crypto -lkrb5support -lcom_err -ldl -lresolv -lpthread
#include "import.h"
*/
import "C"
import (
	"fmt"
	"os"

	"github.com/tarusov/gosmb2/model"
)

// file is hanlder for smb2fh type in Go. Implements File interface.
type file struct {
	path string
	ctx  *context
	ptr  C.fileHandlerPtr
	pos  uint
}

// mkFileHandler creates new smb2fh instance.
func mkFileHandler(ctx *context, path string, mode int) (*file, error) {
	if !ctx.ok() {
		return nil, fmt.Errorf("failed to open file: %v", model.ErrContextIsNil)
	}
	result := C.smb2_open(ctx.ptr, C.CString(path), C.int(mode))
	if result == nil {
		return nil, fmt.Errorf("failed to open file: %v", ctx.lastError())
	}
	return &file{
		path: path,
		ptr:  result,
		ctx:  ctx,
	}, nil
}

// ok check ptr state for instance.
func (f *file) ok() bool {
	return f != nil && f.ptr != nil && f.ctx != nil && f.ctx.ptr != nil
}

// close func free current smb2fh data.
func (f *file) close() {
	if f.ok() {
		C.smb2_close(f.ctx.ptr, f.ptr)
	}
}

// Close impl File interface method.
func (f *file) Close() error {
	f.close()
	return nil
}

// Stat impl File interface method.
func (f *file) Stat() (os.FileInfo, error) {
	if !f.ok() {
		return nil, fmt.Errorf("failed to get file stat: %v", model.ErrContextIsNil)
	}

	return stat(f.ctx, f.path)
}
