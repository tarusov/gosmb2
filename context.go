package gosmb2

/*
#cgo CFLAGS: -I./include
#cgo LDFLAGS: -L./lib -lsmb2
#include "gosmb2.h"
*/
import "C"
import "errors"

// contextPtr is smb2 context type in Go.
type (
	contextPtr C.contextPtr
)

// createContext creates new context instance.
func createContext() (contextPtr, error) {
	result := C.smb2_init_context()
	if result == nil {
		return nil, errors.New("failed to init smb2 context")
	}

	return result, nil
}

// destroyContext free target context.
func destroyContext(ctx contextPtr) {
	if ctx != nil {
		C.smb2_destroy_context(ctx)
	}
}
