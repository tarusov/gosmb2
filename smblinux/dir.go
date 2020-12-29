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

// dir is hanlder for smb2dir type in Go. Implements Dir interface.
type dir struct {
	path string
	ctx  *context
	ptr  C.dirHandlerPtr
	pos  uint
}

// mkDirHandler creates new smb2fh instance.
func mkDirHandler(ctx *context, path string) (*dir, error) {
	if !ctx.ok() {
		return nil, fmt.Errorf("failed to open dir: %v", ErrContextIsNil)
	}
	if path == "." {
		path = ""
	}
	result := C.smb2_opendir(ctx.ptr, C.CString(path))
	if result == nil {
		return nil, fmt.Errorf("failed to open file: %v", lastError(ctx))
	}
	return &dir{
		path: path,
		ptr:  result,
		ctx:  ctx,
	}, nil
}

// ok check ptr state for instance.
func (d *dir) ok() bool {
	return d != nil && d.ptr != nil && d.ctx != nil && d.ctx.ptr != nil
}

// close func free current smb2dir data.
func (d *dir) close() {
	if d.ok() {
		C.smb2_closedir(d.ctx.ptr, d.ptr)
	}
}

// Close impl File interface method.
func (d *dir) Close() error {
	d.close()
	return nil
}

// Stat impl File interface method.
func (d *dir) Stat() (os.FileInfo, error) {
	if !d.ok() {
		return nil, fmt.Errorf("failed to get dir stat: %v", ErrContextIsNil)
	}

	return stat(d.ctx, d.path)
}

// List impl Dir interface method.
func (d *dir) List() ([]*model.DirEntry, error) {
	if !d.ok() {
		return nil, fmt.Errorf("failed to read dir: %v", ErrContextIsNil)
	}

	entries := make([]*model.DirEntry, 0)
	for {
		entry := C.smb2_readdir(d.ctx.ptr, d.ptr)
		if entry == nil {
			return entries, nil
		}

		entryType := model.DirEntryTypeUnknown
		switch entry.st.smb2_type {
		case C.SMB2_TYPE_LINK:
			entryType = model.DirEntryTypeLink
		case C.SMB2_TYPE_DIRECTORY:
			entryType = model.DirEntryTypeDir
		case C.SMB2_TYPE_FILE:
			entryType = model.DirEntryTypeFile
		}

		entries = append(entries, &model.DirEntry{
			Name: C.GoString(entry.name),
			Type: entryType,
			Size: uint64(entry.st.smb2_size),
		})
	}
}
