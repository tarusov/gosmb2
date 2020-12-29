// +build windows

package smbwinnt

import (
	"errors"
	"fmt"
	"syscall"

	"github.com/tarusov/gosmb2/model"
)

// mkContext creates new smb2_context instance.
func mkContext() (uintptr, error) {
	smb2, err := syscall.LoadDLL("smb2.dll")
	if err != nil {
		return 0, fmt.Errorf("failed load smb2.dll: %v", err)
	}

	fmt.Println(smb2)

	proc, err := smb2.FindProc("smb2_init_context")
	if err != nil {
		return 0, fmt.Errorf("failed to find smb2_init_context: %v", err)
	}

	ptr, _, err := proc.Call()
	if err != nil {
		return 0, fmt.Errorf("failed to init smb2 context: %v", err)
	}

	return ptr, nil
}

// Dial create new session with share. URL must be like "//127.0.0.1/share"
func Dial(urlstr string, auth *model.Auth) (model.Share, error) {
	_, err := mkContext()
	if err != nil {
		return nil, fmt.Errorf("failed load dial: %v", err)
	}

	return nil, errors.New("unimplemented yet")
}
