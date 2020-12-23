package smb

/*
#cgo CFLAGS:  -I./include
#cgo amd64   LDFLAGS: -L./lib/amd64 -lsmb2 -lgssapi_krb5 -lkrb5 -lk5crypto -lkrb5support -lcom_err -ldl -lresolv -lpthread
#cgo 386     LDFLAGS: -L./lib/386   -lsmb2 -lgssapi_krb5 -lkrb5 -lk5crypto -lkrb5support -lcom_err -ldl -lresolv -lpthread
#include "import.h"
*/
import "C"
import (
	"fmt"
)

// dir is hanlder for smb2dir type in Go.
type dir struct {
	ptr C.dirHandlerPtr
	ctx *context
}

// mkFile creates new smb2dir instance.
func mkDirHandler(ctx *context, path string) (*dir, error) {
	if !ctx.ok() {
		return nil, ErrContextIsNil
	}
	result := C.smb2_opendir(ctx.ptr, C.CString(path))
	if result == nil {
		return nil, fmt.Errorf("failed to open file: %v", lastError(ctx))
	}
	return &dir{
		ptr: result,
		ctx: ctx,
	}, nil
}

// ok check ptr state for instance.
func (d *dir) ok() bool {
	return d != nil && d.ptr != nil
}

// close func free current smb2dir data.
func (d *dir) close() {
	if d.ok() {
		C.smb2_closedir(d.ctx.ptr, d.ptr)
	}
}
