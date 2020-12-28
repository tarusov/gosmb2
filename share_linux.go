// +build linux

package smb

import (
	"errors"

	"github.com/tarusov/gosmb2/model"
	"github.com/tarusov/gosmb2/smblinux"
)

// Dial create new session with share. URL must be like "//127.0.0.1/share"
func Dial(urlstr string, auth *model.Auth) (model.Share, error) {
	if len(urlstr) == 0 {
		return nil, errors.New("smb url string is empty")
	}

	// Create public auth params.
	if auth == nil {
		auth = &model.Auth{
			Type: model.AuthTypeNegotiate,
		}
	}

	return smblinux.Dial(urlstr, auth)
}
