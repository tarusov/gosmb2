// +build windows
// +build amd64

package smbwinnt

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/tarusov/gosmb2/model"
)

// Read impl File interface method.
//
// int smb2_pread(struct smb2_context *smb2, struct smb2fh *fh,
// 	uint8_t *buf, uint32_t count, uint64_t offset);
func (f *file) Read(p []byte) (n int, err error) {
	if !f.ok() {
		return 0, fmt.Errorf("failed to read file: %v", model.ErrContextIsNil)
	}

	bufSizeProc, err := f.ctx.dll.FindProc("smb2_get_max_read_size")
	if err != nil {
		return 0, fmt.Errorf("failed to read file: %v", err)
	}
	readProc, err := f.ctx.dll.FindProc("smb2_pread")
	if err != nil {
		return 0, fmt.Errorf("failed to read file: %v", err)
	}

	maxBufSize, _, err := bufSizeProc.Call(f.ctx.ptr)
	if err != syscall.Errno(0) {
		return 0, fmt.Errorf("failed to read file: %v", err)
	}

	if pSize := len(p); pSize < int(maxBufSize) {
		maxBufSize = uintptr(pSize)
	}

	var (
		bufChunk    = make([]byte, maxBufSize)
		bufChunkPtr = uintptr(unsafe.Pointer(&bufChunk[0]))
	)

	for {
		count, _, err := readProc.Call(
			f.ctx.ptr,
			f.ptr,
			bufChunkPtr,
			maxBufSize,
			f.pos,
		)

		if err != syscall.Errno(0) {
			return 0, fmt.Errorf("failed to read file: %v", err)
		}
		if count == 0 {
			return n, nil // finished successful. EOF.
		}
		if count < 0 {
			return 0, f.ctx.lastError() // recv error.
		}

		// Copy to p from chunk.
		for i := 0; i < int(count); i++ {
			p[i+n] = bufChunk[i]
		}

		n += int(count) // inc read bytes.
		f.pos += count  // move file pos.
	}
}

// Stat impl File interface method.
//
/*
 * smb2_seek() SEEK_SET and SEEK_CUR are fully supported.
 * SEEK_END only returns the end-of-file from the original open.
 * (it will not call fstat to discover the current file size and will not block)
 */
//  int64_t smb2_lseek(struct smb2_context *smb2, struct smb2fh *fh,
// 	int64_t offset, int whence, uint64_t *current_offset);
func (f *file) Seek(offset int64, whence int) (int64, error) {

	fmt.Println("!!! SEEK")

	if !f.ok() {
		return 0, fmt.Errorf("failed to seek file pos: %v", model.ErrContextIsNil)
	}

	seekProc, err := f.ctx.dll.FindProc("smb2_lseek")
	if err != nil {
		return 0, fmt.Errorf("failed to seek file pos: %v", err)
	}

	shifted, _, err := seekProc.Call(
		f.ctx.ptr,
		f.ptr,
		uintptr(offset),
		uintptr(whence),
		f.pos,
	)
	if err != syscall.Errno(0) {
		return 0, fmt.Errorf("failed to read file: %v", err)
	}

	f.pos = shifted

	return int64(f.pos), nil
}
