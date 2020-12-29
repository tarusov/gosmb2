package model

// Error is lightweight error type.
type Error string

// Error impl error interface.
func (e Error) Error() string {
	return string(e)
}

// Error constants.
const (
	ErrContextIsNil        Error = "smb context is nilptr"
	ErrUnknownAuthType     Error = "unknown auth type"
	ErrInvalidResourcePath Error = "invalid resource path"
)
