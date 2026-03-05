package api

import (
	"encoding/json"
	"net/http"

	"github.com/dcdavidev/bastion/packages/auth"
	"github.com/dcdavidev/bastion/packages/models"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
)

// PasskeyRegisterBegin generates registration options for a new passkey.
func (h *Handler) PasskeyRegisterBegin(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(auth.UserKey).(uuid.UUID)

	user, err := h.DB.GetUserByID(r.Context(), uid)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Fetch existing credentials
	creds, err := h.DB.GetWebAuthnCredentials(r.Context(), user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	webauthUser := &WebAuthnUser{User: user, Credentials: creds}

	options, session, err := h.WebAuthn.BeginRegistration(webauthUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Store session data using user ID as key
	h.sessions.Store("reg_"+user.ID.String(), session)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(options)
}

// PasskeyRegisterFinish finalizes passkey registration.
func (h *Handler) PasskeyRegisterFinish(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(auth.UserKey).(uuid.UUID)

	user, err := h.DB.GetUserByID(r.Context(), uid)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Retrieve session data
	sessionData, ok := h.sessions.Load("reg_" + user.ID.String())
	if !ok {
		http.Error(w, "Registration session not found", http.StatusBadRequest)
		return
	}
	session := sessionData.(*webauthn.SessionData)
	h.sessions.Delete("reg_" + user.ID.String())

	webauthUser := &WebAuthnUser{User: user}

	credential, err := h.WebAuthn.FinishRegistration(webauthUser, *session, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save credential to DB
	webauthCred := &models.WebAuthnCredential{
		ID:              credential.ID,
		PublicKey:       credential.PublicKey,
		AttestationType: credential.AttestationType,
		Transport:       h.fromWebAuthnTransport(credential.Transport),
		SignCount:       credential.Authenticator.SignCount,
		CloneWarning:    credential.Authenticator.CloneWarning,
	}

	err = h.DB.AddWebAuthnCredential(r.Context(), user.ID, webauthCred)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// PasskeyLoginBegin generates authentication options for a passkey login.
func (h *Handler) PasskeyLoginBegin(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	user, _, _, err := h.DB.GetUserByEmail(r.Context(), email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	creds, err := h.DB.GetWebAuthnCredentials(r.Context(), user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	webauthUser := &WebAuthnUser{User: user, Credentials: creds}

	options, session, err := h.WebAuthn.BeginLogin(webauthUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Store session data using user ID as key
	h.sessions.Store("login_"+user.ID.String(), session)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(options)
}

// PasskeyLoginFinish finalizes passkey authentication and returns a JWT.
func (h *Handler) PasskeyLoginFinish(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	user, _, _, err := h.DB.GetUserByEmail(r.Context(), email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Retrieve session data
	sessionData, ok := h.sessions.Load("login_" + user.ID.String())
	if !ok {
		http.Error(w, "Login session not found", http.StatusBadRequest)
		return
	}
	session := sessionData.(*webauthn.SessionData)
	h.sessions.Delete("login_" + user.ID.String())

	creds, _ := h.DB.GetWebAuthnCredentials(r.Context(), user.ID)
	webauthUser := &WebAuthnUser{User: user, Credentials: creds}

	credential, err := h.WebAuthn.FinishLogin(webauthUser, *session, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Update sign count
	for i := range creds {
		if string(creds[i].ID) == string(credential.ID) {
			creds[i].SignCount = credential.Authenticator.SignCount
			h.DB.UpdateWebAuthnCredential(r.Context(), &creds[i])
			break
		}
	}

	// Generate JWT
	token, err := auth.GenerateToken(user.ID, user.Username, user.Role == "ADMIN")
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *Handler) fromWebAuthnTransport(t []protocol.AuthenticatorTransport) []string {
	res := make([]string, len(t))
	for i, transport := range t {
		res[i] = string(transport)
	}
	return res
}
