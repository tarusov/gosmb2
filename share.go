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
	"os"
	"strings"
)

// Share interface is handler for connection with share.
type Share interface {
	OpenFile(filename string, mode int) (File, error)

	Close() // Close established connection.
}

// Dial create new session with share. URL must be like "//127.0.0.1/qqck"
func Dial(urlstr string, auth *Auth) (Share, error) {
	ctx, err := mkContext()
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}

	return mkShare(ctx, auth, urlstr)
}

// conn is connection state handler.
type share struct {
	ctx  *context
	auth *Auth
}

// mkConn create new connection instance.
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

	C.smb2_set_domain(ctx.ptr, C.CString(auth.Domain))
	C.smb2_set_password(ctx.ptr, C.CString(auth.Password))
	if host, err := os.Hostname(); err == nil {
		C.smb2_set_workstation(ctx.ptr, C.CString(host))
	}

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

	if result := C.smb2_connect_share(
		ctx.ptr,
		C.CString(uri.Host),
		C.CString(parts[0]),
		C.CString(auth.Username),
	); result < 0 {
		return nil, fmt.Errorf("failed to connect share: %v", lastError(ctx))
	}
	return &share{
		ctx:  ctx,
		auth: auth,
	}, nil
}

// ok check ptr state for instance.
func (s *share) ok() bool {
	return s != nil && s.ctx != nil
}

// close close current connection.
func (s *share) disconnect() {
	if s.ok() {
		C.smb2_disconnect_share(s.ctx.ptr)
	}
}

// Close impl Conn interface,
// free all resources for connection.
func (s *share) Close() {
	s.disconnect()
	s.ctx.free() // keep last!
}

func (s *share) OpenFile(filename string, mode int) (File, error) {
	if !s.ok() {
		return nil, fmt.Errorf("failed to open file: %v", ErrContextIsNil)
	}

	return mkFileHandler(s.ctx, filename, mode)
}
