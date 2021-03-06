// +build linux
// +build amd64

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
	"io"
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

	maxBufSize := C.smb2_get_max_read_size(f.ctx.ptr)
	if pSize := len(p); pSize < int(maxBufSize) {
		maxBufSize = C.uint(pSize)
	}

	var (
		bufChunk    = make([]byte, maxBufSize)
		bufChunkPtr = (*C.uchar)(&bufChunk[0])
	)

	for {
		count := C.smb2_pread(
			f.ctx.ptr,
			f.ptr,
			bufChunkPtr,
			maxBufSize,
			C.ulong(f.pos),
		)
		if count == 0 {
			return n, nil // finished successful.
		}
		if count == -(C.EAGAIN) {
			continue // need to read again.
		}
		if count < 0 {
			return 0, f.ctx.lastError() // recv error.
		}

		// Copy to p from chunk.
		for i := 0; i < int(count); i++ {
			p[i+n] = bufChunk[i]
		}

		n += int(count)      // inc read bytes.
		f.pos += uint(count) // move file pos.
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
	fwh := C.SEEK_SET
	switch whence {
	case io.SeekStart:
		// Go on, already set.
	case io.SeekCurrent:
		fwh = C.SEEK_CUR
	case io.SeekEnd:
		fwh = C.SEEK_END
	}

	pos := C.smb2_lseek(
		f.ctx.ptr,
		f.ptr,
		C.int64_t(offset),
		C.int(fwh),
		C.ulong(f.pos),
	)
	if pos < 0 {
		return 0, fmt.Errorf("failed to seek pos into file: %v", f.ctx.lastError())
	}

	f.pos = uint(pos) // shift offset.

	return int64(pos), nil
}

// conn create new connection with server.
func conn(ctx *context, host, share, user string) error {
	if !ctx.ok() {
		return nil, fmt.Errorf("failed to connect: %v", model.ErrContextIsNil)
	}

	hostPtr, err := syscall.BytePtrFromString(host)
	if err != nil {
		return 0, err
	}
	sharePtr, err := syscall.BytePtrFromString(share)
	if err != nil {
		return 0, err
	}
	userPtr, err := syscall.BytePtrFromString(user)
	if err != nil {
		return 0, err
	}

	if result := C.smb2_connect_share(
		ctx.ptr,
		uintptr(unsafe.Pointer(hostPtr)),
		uintptr(unsafe.Pointer(sharePtr)),
		uintptr(unsafe.Pointer(userPtr)),
	); result < 0 {
		return fmt.Errorf("failed to connect share: %d %v", result, ctx.lastError())
	}

	return nil
}
