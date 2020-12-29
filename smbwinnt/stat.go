// +build windows

package smbwinnt

import (
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"

	"github.com/tarusov/gosmb2/model"
)

// Avaliable samba entry types.
const (
	EntryTypeFile = iota
	EntryTypeDir
	EntryTypeUnknown
)

/*
struct smb2_stat_64 {
    uint32_t smb2_type;
    uint32_t smb2_nlink;
    uint64_t smb2_ino;
    uint64_t smb2_size;
	uint64_t smb2_atime;
	uint64_t smb2_atime_nsec;
	uint64_t smb2_mtime;
	uint64_t smb2_mtime_nsec;
	uint64_t smb2_ctime;
	uint64_t smb2_ctime_nsec;
    uint64_t smb2_btime;
    uint64_t smb2_btime_nsec;
};

*/

// Get stat about file.
func stat(ctx *context, path string) (os.FileInfo, error) {
	if !ctx.ok() {
		return nil, fmt.Errorf("failed to get stat: %v", model.ErrContextIsNil)
	}
	proc, err := ctx.dll.FindProc("smb2_stat")
	if err != nil {
		return nil, err
	}
	pathPtr, err := syscall.BytePtrFromString(path)
	if err != nil {
		return nil, err
	}

	// See info below.
	data := &struct {
		Type        uint32
		_           uint32
		_           uint64
		Size        uint64
		_           uint64
		_           uint64
		ModTime     uint64
		ModTimeNsec uint64
		_           uint64
		_           uint64
		_           uint64
		_           uint64
	}{}

	_, _, err = proc.Call(
		ctx.ptr,
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(unsafe.Pointer(data)),
	)
	if err != syscall.Errno(0) {
		return nil, err
	}

	var entryType int
	switch data.Type {
	case 0:
		entryType = EntryTypeFile
	case 1:
		entryType = EntryTypeDir
	default:
		entryType = EntryTypeUnknown
	}

	return &entryInfo{
		name:      path,
		entryType: entryType,
		size:      int64(data.Size),
		modTime:   time.Unix(int64(data.ModTime), int64(data.ModTimeNsec)),
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
