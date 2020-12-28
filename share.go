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
	"net/url"
	"strings"
)

// Share interface is handler for connection with share.
type Share interface {
	Echo() error                                      // Send echo request to server.
	OpenFile(filename string, mode int) (File, error) // Open shared file.
	Close()                                           // Close established connection.
}

// Dial create new session with share. URL must be like "//127.0.0.1/share"
func Dial(urlstr string, auth *Auth) (Share, error) {

	fmt.Println("!!! Dial")

	ctx, err := mkContext()
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}

	return mkShare(ctx, auth, urlstr)
}

// share is connection state handler.
type share struct {
	ctx *context
}

// mkShare create new connection instance.
func mkShare(ctx *context, auth *Auth, urlstr string) (*share, error) {
	if !ctx.ok() {
		return nil, fmt.Errorf("failed to connect share: %v", ErrContextIsNil)
	}
	if auth == nil {
		return nil, fmt.Errorf("failed to connect share: credential not provided")
	}

	switch auth.Type {
	case AuthTypeKrb5:
		C.smb2_set_authentication(ctx.ptr, C.SMB2_SEC_KRB5)
		C.smb2_set_security_mode(ctx.ptr, C.SMB2_NEGOTIATE_SIGNING_REQUIRED)
	case AuthTypeNTLM:
		C.smb2_set_authentication(ctx.ptr, C.SMB2_SEC_NTLMSSP)
		C.smb2_set_security_mode(ctx.ptr, C.SMB2_NEGOTIATE_SIGNING_REQUIRED)
	default:
		C.smb2_set_security_mode(ctx.ptr, C.SMB2_NEGOTIATE_SIGNING_ENABLED)
	}

	if len(auth.Domain) > 0 {
		C.smb2_set_domain(ctx.ptr, C.CString(auth.Domain))
	}
	if len(auth.Password) > 0 {
		C.smb2_set_password(ctx.ptr, C.CString(auth.Password))
	}
	/*
		if host, err := os.Hostname(); err == nil {
			C.smb2_set_workstation(ctx.ptr, C.CString(host))
		}*/

	uri, err := url.Parse(urlstr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect share: %v", err)
	}

	// extract share name from uri ()
	parts := strings.Split(uri.Path, "/")
	if len(parts) == 0 {
		return nil, fmt.Errorf("failed to connect share: %v", ErrInvalidResourcePath)
	}
	if len(parts[0]) == 0 {
		parts = parts[1:]
	}

	fmt.Println("!!! connect to", parts[0])

	if result := C.smb2_connect_share(
		ctx.ptr,
		C.CString(uri.Host),
		C.CString(parts[0]),
		C.CString(auth.Username),
	); result < 0 {
		return nil, fmt.Errorf("failed to connect share: %d %v", result, lastError(ctx))
	}

	return &share{
		ctx: ctx,
	}, nil
}

// ok check ptr state for instance.
func (s *share) ok() bool {
	return s != nil && s.ctx != nil
}

// disconnect close current connection.
func (s *share) disconnect() {
	if s.ok() {
		C.smb2_disconnect_share(s.ctx.ptr)
	}
}

// Close impl Share interface method,
// free all resources for connection.
func (s *share) Close() {
	s.disconnect()
	s.ctx.free() // keep last!
}

// OpenFile impl Share interface method.
func (s *share) OpenFile(filename string, mode int) (File, error) {
	if !s.ok() {
		return nil, fmt.Errorf("failed to open file: %v", ErrContextIsNil)
	}

	return mkFileHandler(s.ctx, filename, mode)
}

// Echo send request.
func (s *share) Echo() error {
	if !s.ok() {
		return fmt.Errorf("failed to send echo request: %v", ErrContextIsNil)
	}

	fmt.Println("!!! Echo")

	result := C.smb2_echo(s.ctx.ptr)
	if result != 0 {
		return fmt.Errorf("failed to send echo request: %d %v", result, lastError(s.ctx))
	}

	return nil
}
