package smblinux

/*
#cgo CFLAGS:  -I./include
#cgo amd64   LDFLAGS: -L./lib/amd64 -lsmb2 -lgssapi_krb5 -lkrb5 -lk5crypto -lkrb5support -lcom_err -ldl -lresolv -lpthread
#cgo 386     LDFLAGS: -L./lib/386   -lsmb2 -lgssapi_krb5 -lkrb5 -lk5crypto -lkrb5support -lcom_err -ldl -lresolv -lpthread
#include "import.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/tarusov/gosmb2/model"
)

// context is handler for smb2_context type in Go.
type context struct {
	ptr C.contextPtr
}

// mkContext creates new smb2_context instance.
func mkContext() (*context, error) {
	result := C.smb2_init_context()
	if result == nil {
		return nil, errors.New("failed to init smb2 context")
	}

	return &context{
		ptr: result,
	}, nil
}

// ok check ptr state for instance.
func (c *context) ok() bool {
	return c != nil && c.ptr != nil
}

// Free current smb2_context.
func (c *context) free() {
	if c.ok() {
		C.smb2_destroy_context(c.ptr)
	}
	c.ptr = nil
}

// share is connection state handler.
type share struct {
	ctx *context
}

// Dial create new session with share. URL must be like "//127.0.0.1/share"
func Dial(urlstr string, auth *model.Auth) (model.Share, error) {
	ctx, err := mkContext()
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}

	switch auth.Type {
	case model.AuthTypeKrb5:
		C.smb2_set_authentication(ctx.ptr, C.SMB2_SEC_KRB5)
		C.smb2_set_security_mode(ctx.ptr, C.SMB2_NEGOTIATE_SIGNING_REQUIRED)
	case model.AuthTypeNTLM:
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

	err = mkConnect(ctx, uri.Host, parts[0], auth.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
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
func (s *share) Close() error {
	s.disconnect()
	s.ctx.free() // keep last!
	return nil
}

// OpenFile impl Share interface method.
func (s *share) OpenFile(filename string, mode int) (model.File, error) {
	if !s.ok() {
		return nil, fmt.Errorf("failed to open file: %v", ErrContextIsNil)
	}

	return mkFileHandler(s.ctx, filename, mode)
}

// OpenDir impl Share interface method.
func (s *share) OpenDir(filename string) (model.Dir, error) {
	if !s.ok() {
		return nil, fmt.Errorf("failed to open file: %v", ErrContextIsNil)
	}

	return mkDirHandler(s.ctx, filename)
}

// Echo send request.
func (s *share) Echo() error {
	if !s.ok() {
		return fmt.Errorf("failed to send echo request: %v", ErrContextIsNil)
	}

	result := C.smb2_echo(s.ctx.ptr)
	if result != 0 {
		return fmt.Errorf("failed to send echo request: %d %v", result, s.ctx.lastError())
	}

	return nil
}
