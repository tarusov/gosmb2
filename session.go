package gosmb2

/*
#cgo CFLAGS: -I./include
#cgo LDFLAGS: -L./lib -lsmb2
#include "gosmb2.h"
*/
import "C"

type smb2Session struct {
	ctx contextPtr
	url urlPtr
}

// NewSession create new sync file-read session for smb.
func NewSession(opts *SessionOptions) (Session, error) {
	ctx, err := createContext()
	if err != nil {
		return nil, err
	}

	url, err := createURL(ctx, opts.dsn())
	if err != nil {
		return nil, err
	}

	err = setAuthType(ctx, opts.AuthType)
	if err != nil {
		return nil, err
	}

	err = setNegotiateSigning(ctx, opts.SecurityMode)
	if err != nil {
		return nil, err
	}

	return &smb2Session{
		ctx: ctx,
		url: url,
	}, nil
}

// Close
func (s *smb2Session) Close() error {
	destroyURL(s.url)
	destroyContext(s.ctx)

	return nil
}
