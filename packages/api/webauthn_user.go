package api

import (
	"github.com/dcdavidev/bastion/packages/models"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

// WebAuthnUser wraps a models.User to implement webauthn.User interface.
type WebAuthnUser struct {
	User        *models.User
	Credentials []models.WebAuthnCredential
}

func (u *WebAuthnUser) WebAuthnID() []byte {
	return u.User.ID[:]
}

func (u *WebAuthnUser) WebAuthnName() string {
	return u.User.Username
}

func (u *WebAuthnUser) WebAuthnDisplayName() string {
	return u.User.Username
}

func (u *WebAuthnUser) WebAuthnIcon() string {
	return ""
}

func (u *WebAuthnUser) WebAuthnCredentials() []webauthn.Credential {
	res := make([]webauthn.Credential, len(u.Credentials))
	for i, c := range u.Credentials {
		res[i] = webauthn.Credential{
			ID:              c.ID,
			PublicKey:       c.PublicKey,
			AttestationType: c.AttestationType,
			Transport:       u.toWebAuthnTransport(c.Transport),
			Authenticator: webauthn.Authenticator{
				SignCount: c.SignCount,
				CloneWarning: c.CloneWarning,
			},
		}
	}
	return res
}

func (u *WebAuthnUser) toWebAuthnTransport(t []string) []protocol.AuthenticatorTransport {
	res := make([]protocol.AuthenticatorTransport, len(t))
	for i, transport := range t {
		res[i] = protocol.AuthenticatorTransport(transport)
	}
	return res
}
