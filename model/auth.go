package model

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
