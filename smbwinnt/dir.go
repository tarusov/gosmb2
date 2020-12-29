// +build windows

package smbwinnt

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/tarusov/gosmb2/model"
)

// dir is hanlder for smb2dir type in Go. Implements Dir interface.
type dir struct {
	path string
	ctx  *context
	ptr  uintptr
}

// mkDirHandler creates new smb2fh instance.
func mkDirHandler(ctx *context, path string) (*dir, error) {
	if !ctx.ok() {
		return nil, fmt.Errorf("failed to open dir: %v", model.ErrContextIsNil)
	}

	proc, err := ctx.dll.FindProc("smb2_opendir")
	if err != nil {
		return nil, err
	}
	if path == "." {
		path = ""
	}
	pathPtr, err := syscall.BytePtrFromString(path)
	if err != nil {
		return nil, err
	}

	ptr, _, err := proc.Call(
		ctx.ptr,
		uintptr(unsafe.Pointer(pathPtr)),
	)
	if ptr == 0 || err != syscall.Errno(0) {
		return nil, ctx.lastError()
	}

	return &dir{
		path: path,
		ptr:  ptr,
		ctx:  ctx,
	}, nil
}

// ok check ptr state for instance.
func (d *dir) ok() bool {
	return d != nil && d.ptr != 0 && d.ctx != nil && d.ctx.ptr != 0
}

// close func free current smb2dir data.
func (d *dir) close() error {
	if d.ok() {
		proc, err := d.ctx.dll.FindProc("smb2_closedir")
		if err != nil {
			return err
		}

		_, _, err = proc.Call(d.ctx.ptr, d.ptr)
		if err != syscall.Errno(0) {
			return err
		}
	}

	return nil
}

// Close impl Dir interface method.
func (d *dir) Close() error {
	return d.close()
}

// Stat impl Dir interface method.
func (d *dir) Stat() (os.FileInfo, error) {
	if !d.ok() {
		return nil, fmt.Errorf("failed to get dir stat: %v", model.ErrContextIsNil)
	}

	return stat(d.ctx, d.path)
}

// List impl Dir interface method.
func (d *dir) List() ([]*model.DirEntry, error) {
	if !d.ok() {
		return nil, fmt.Errorf("failed to read dir: %v", model.ErrContextIsNil)
	}

	entries := make([]*model.DirEntry, 0)
	return entries, nil

	/*
		for {
			// TODO: get ptr to struct from proc.Call


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
	*/
}
