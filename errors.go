package smb

/*
#cgo CFLAGS:  -I./include
#cgo amd64   LDFLAGS: -L./lib/amd64 -lsmb2 -lgssapi_krb5 -lkrb5 -lk5crypto -lkrb5support -lcom_err -ldl -lresolv -lpthread
#cgo 386     LDFLAGS: -L./lib/386   -lsmb2 -lgssapi_krb5 -lkrb5 -lk5crypto -lkrb5support -lcom_err -ldl -lresolv -lpthread
#include "import.h"
*/
import "C"
import "errors"

// Error is lightweight error type.
type Error string

// Error impl error interface.
func (e Error) Error() string {
	return string(e)
}

// Error constants.
const (
	ErrContextIsNil        Error = "smb context is nilptr"
	ErrUnknownAuthType     Error = "unknown auth type"
	ErrInvalidResourcePath Error = "invalid resource path"
)

// lastError returns last smb2 error from smb2_context in go format.
func lastError(ctx *context) error {
	result := C.smb2_get_error(ctx.ptr)
	return errors.New(C.GoString(result))
}
