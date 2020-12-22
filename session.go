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

func (s smb2Session) Connect() error {
	result := C.smb2_connect_share(
		s.ctx,
		(*C.char)(s.url.server),
		(*C.char)(s.url.share),
		(*C.char)(s.url.user),
	)
	if result != 0 {
		return lastError(s.ctx)
	}

	fh := C.smb2_open(s.ctx, s.url.path, C.O_RDONLY)
	if fh == nil {
		return lastError(s.ctx)
	}
	/*
	   while ((count = smb2_pread(smb2, fh, buf, MAXBUF, pos)) != 0) {
	           if (count == -EAGAIN) {
	                   continue;
	           }
	           if (count < 0) {
	                   fprintf(stderr, "Failed to read file. %s\n",
	                           smb2_get_error(smb2));
	                   rc = 1;
	                   break;
	           }
	           write(0, buf, count);
	           pos += count;
	   };
	*/

	defer func() {
		C.smb2_close(s.ctx, fh)
		C.smb2_disconnect_share(s.ctx)
	}()

	return nil
}

// Close current session.
func (s *smb2Session) Close() error {
	destroyURL(s.url)
	destroyContext(s.ctx)

	return nil
}
