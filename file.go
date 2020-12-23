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
	"io"
	"os"
)

// File is interface for fileHandler.
type File interface {
	// Stat returns file data.
	Stat() (os.FileInfo, error)

	// Close current file.
	Close()
}

// file is hanlder for smb2fh type in Go.
type file struct {
	path string
	ctx  *context
	ptr  C.fileHandlerPtr
	pos  C.ulonglong
}

// mkFileHandler creates new smb2fh instance.
func mkFileHandler(ctx *context, path string, mode int) (*file, error) {
	if !ctx.ok() {
		return nil, fmt.Errorf("failed to open file: %v", ErrContextIsNil)
	}
	result := C.smb2_open(ctx.ptr, C.CString(path), C.int(mode))
	if result == nil {
		return nil, fmt.Errorf("failed to open file: %v", lastError(ctx))
	}
	return &file{
		ptr: result,
		ctx: ctx,
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
func (f *file) Close() {
	f.close()
}

// Stat impl File interface method.
func (f *file) Stat() (os.FileInfo, error) {
	if !f.ok() {
		return nil, fmt.Errorf("failed to get file stat: %v", ErrContextIsNil)
	}

	return stat(f.ctx, f.path)
}

// Read impl File interface method.
//
// int smb2_pread(struct smb2_context *smb2, struct smb2fh *fh,
// 	uint8_t *buf, uint32_t count, uint64_t offset);
func (f *file) Read(p []byte) (n int, err error) {
	if !f.ok() {
		return 0, fmt.Errorf("failed to read file: %v", ErrContextIsNil)
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
			f.pos,
		)
		if count == 0 {
			return n, nil // finished successful.
		}
		if count == -(C.EAGAIN) {
			continue // need to read again.
		}
		if count < 0 {
			return 0, lastError(f.ctx) // recv error.
		}

		// Copy to p from chunk.
		for i := 0; i < len(bufChunk); i++ {
			p[i+n] = bufChunk[i]
		}

		n += int(count)             // inc read bytes.
		f.pos += C.ulonglong(count) // move file pos.
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
		C.longlong(offset),
		C.int(fwh),
		&f.pos)
	if pos < 0 {
		return 0, fmt.Errorf("failed to seek pos into file: %v", lastError(f.ctx))
	}

	f.pos = C.ulonglong(pos) // shift offset.

	return int64(pos), nil
}
