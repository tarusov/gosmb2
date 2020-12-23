package smb

/*
#cgo CFLAGS:  -I./include
#cgo amd64   LDFLAGS: -L./lib/amd64 -lsmb2 -lgssapi_krb5 -lkrb5 -lk5crypto -lkrb5support -lcom_err -ldl -lresolv -lpthread
#cgo 386     LDFLAGS: -L./lib/386   -lsmb2 -lgssapi_krb5 -lkrb5 -lk5crypto -lkrb5support -lcom_err -ldl -lresolv -lpthread
#include "import.h"
*/
import "C"

// AuthType is supported auth type.
type AuthType string

// Available auth types.
const (
	AuthTypeNegotiate AuthType = "negotiate" // server-side auth type.
	AuthTypeNTLM      AuthType = "ntlm"      // basic ntlm auth.
	AuthTypeKrb5      AuthType = "krb5"      // kerberos from kinit.
)

// Auth is struct for auth paramets. Use func for create
type Auth struct {
	Type     AuthType
	Domain   string
	Username string
	Password string
}
