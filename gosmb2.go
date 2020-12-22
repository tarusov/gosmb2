package gosmb2

/*
#cgo CFLAGS: -I./include
#cgo LDFLAGS: -L./lib -lsmb2
#include "gosmb2.h"
*/
import "C"

// Session interface declares smb session methods.
type Session interface {
	// Connect launch connection for session.
	Connect() error

	// Close connection and free context data.
	Close() error
}
