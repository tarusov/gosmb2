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
	"time"
)

// Avaliable samba entry types.
const (
	EntryTypeFile = iota
	EntryTypeDir
	EntryTypeUnknown
)

func stat(ctx *context, path string) (os.FileInfo, error) {
	if !ctx.ok() {
		return nil, fmt.Errorf("failed to get stat: %v", ErrContextIsNil)
	}

	var (
		data C.fileInfo
		//vfs  C.vfsInfo
	)

	result := C.smb2_stat(ctx.ptr, C.CString(path), &data)
	if result < 0 {
		return nil, fmt.Errorf("failed to get stat: %d %v", result, lastError(ctx))
	}

	var entryType int
	switch data.smb2_type {
	case C.SMB2_TYPE_FILE:
		entryType = EntryTypeFile
	case C.SMB2_TYPE_DIRECTORY:
		entryType = EntryTypeDir
	default:
		entryType = EntryTypeUnknown
	}

	return &entryInfo{
		name:      path,
		entryType: entryType,
		size:      int64(data.smb2_size),
		modTime:   time.Unix(int64(data.smb2_mtime), 0),
	}, nil
}

type entryInfo struct {
	size      int64
	name      string
	entryType int
	modTime   time.Time
}

func (i *entryInfo) Name() string {
	return i.name
}

func (i *entryInfo) Size() int64 {
	return i.size
}

func (i *entryInfo) Mode() os.FileMode {
	return 0
}

func (i *entryInfo) ModTime() time.Time {
	return i.modTime
}

func (i *entryInfo) IsDir() bool {
	return i.entryType == EntryTypeDir
}

func (i *entryInfo) Sys() interface{} {
	return nil
}
