package gosmb2

/*
#cgo CFLAGS: -I./include
#cgo LDFLAGS: -L./lib -lsmb2
#include "gosmb2.h"
*/
import "C"
import (
	"errors"
	"fmt"
)

// SessionOptions defines SMB2 session options.
type SessionOptions struct {
	Server       string
	Port         int
	Domain       string
	User         string
	Path         string
	AuthType     AuthType
	SecurityMode NegotiateSigningType
}

// AuthType is smb2 auth types.
type AuthType string

// Available auth types.
const (
	AuthTypeUndefined AuthType = "undefined"
	AuthTypeNTLM      AuthType = "ntlm"
	AuthTypeKerberos  AuthType = "kerberos"
)

// setAuthType set auth type into smb context.
func setAuthType(ctx contextPtr, auth AuthType) error {
	if ctx == nil {
		return errors.New("smb context is nil")
	}
	/*
		switch auth {
		case AuthTypeUndefined, "":
			ctx.sec = C.SMB2_SEC_UNDEFINED
		case AuthTypeNTLM:
			ctx.sec = C.SMB2_SEC_NTLMSSP
		case AuthTypeKerberos:
			ctx.sec = C.SMB2_SEC_KERBEROS
		default:
			return fmt.Errorf("unknown auth type %q", auth)
		}
	*/
	return nil
}

// NegotiateSigningType is smb2 negotiate type.
type NegotiateSigningType string

// Available negotiate signing types.
const (
	NegotiateSigningTypeEnabled  NegotiateSigningType = "enabled"
	NegotiateSigningTypeRequired NegotiateSigningType = "required"
)

// Set the security mode for the connection.
func setNegotiateSigning(ctx contextPtr, ng NegotiateSigningType) error {
	if ctx == nil {
		return errors.New("smb context is nil")
	}

	switch ng {
	case NegotiateSigningTypeEnabled:
		C.smb2_set_security_mode(ctx, C.SMB2_NEGOTIATE_SIGNING_ENABLED)
	case NegotiateSigningTypeRequired:
		C.smb2_set_security_mode(ctx, C.SMB2_NEGOTIATE_SIGNING_REQUIRED)
	default:
		return fmt.Errorf("unknown security mode %q", ng)
	}

	return nil
}

// dsn create resource url from options.
func (o *SessionOptions) dsn() string {
	var domain string
	if len(o.Domain) > 0 {
		domain = o.Domain + ";"
	}

	var user string
	if len(o.User) > 0 {
		user = o.User + "@"
	}

	var port int
	if o.Port > 0 && o.Port < 0xFFFF {
		port = o.Port
	}

	var share string
	if len(o.Path) > 0 {
		if o.Path[0] != '/' {
			share = share + "/"
		}
		share = o.Path
	} else {
		share = "/"
	}

	return fmt.Sprintf("smb://%s%s%s:%d%s", domain, user, o.Server, port, share)
}
