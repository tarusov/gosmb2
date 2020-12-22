package gosmb2

/*
#cgo CFLAGS: -I./include
#cgo LDFLAGS: -L./lib -lsmb2
#include "gosmb2.h"
*/
import "C"
import (
	"errors"
)

// urlPtr is smb2 url type in Go.
type urlPtr C.urlPtr

// createURL create new smb2
func createURL(ctx contextPtr, url string) (urlPtr, error) {
	if ctx == nil {
		return nil, errors.New("smb context is nil")
	}

	result := C.smb2_parse_url(ctx, C.CString(url))
	if result == nil {
		return nil, lastError(ctx)
	}

	return result, nil
}

func destroyURL(url urlPtr) {
	if url != nil {
		C.smb2_destroy_url(url)
	}
}
