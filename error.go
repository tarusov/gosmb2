package gosmb2

/*
#cgo CFLAGS: -I./include
#cgo LDFLAGS: -L./lib -lsmb2
#include "gosmb2.h"
*/
import "C"
import "errors"

// lastError returns last smb2 error from context in go format.
func lastError(ctx contextPtr) error {
	cstr := C.smb2_get_error(ctx)
	return errors.New(C.GoString(cstr))
}
