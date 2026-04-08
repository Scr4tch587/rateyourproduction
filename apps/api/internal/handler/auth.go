package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/scr4tch/rateyourproduction/apps/api/internal/model"
)

const sessionDuration = 30 * 24 * time.Hour

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var req model.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Username == "" || req.Email == "" || len(req.Password) < 8 {
		writeError(w, http.StatusBadRequest, "username, email, and password are required")
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create account")
		return
	}

	ctx := r.Context()
	tx, err := h.DB.Begin(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to start signup")
		return
	}
	defer tx.Rollback(ctx)

	var isFirstUser bool
	if err := tx.QueryRow(ctx, `SELECT count(*) = 0 FROM profiles`).Scan(&isFirstUser); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create account")
		return
	}

	var profileID string
	err = tx.QueryRow(ctx, `
		INSERT INTO profiles (username, display_name, is_admin)
		VALUES ($1, $2, $3)
		RETURNING id
	`, req.Username, req.DisplayName, isFirstUser).Scan(&profileID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "username is already in use")
		return
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO auth_accounts (profile_id, email, password_hash)
		VALUES ($1, $2, $3)
	`, profileID, req.Email, string(passwordHash))
	if err != nil {
		writeError(w, http.StatusBadRequest, "email is already in use")
		return
	}

	token, err := generateToken()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create session")
		return
	}
	expiresAt := time.Now().Add(sessionDuration)
	_, err = tx.Exec(ctx, `
		INSERT INTO sessions (profile_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`, profileID, hashToken(token), expiresAt)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create session")
		return
	}

	if err := tx.Commit(ctx); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create account")
		return
	}

	setSessionCookie(w, token, expiresAt)
	writeJSON(w, http.StatusCreated, model.SessionProfile{
		ID:          profileID,
		Username:    req.Username,
		Email:       req.Email,
		DisplayName: req.DisplayName,
		IsAdmin:     isFirstUser,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	var user model.SessionProfile
	var passwordHash string
	err := h.DB.QueryRow(r.Context(), `
		SELECT p.id, p.username, a.email, p.display_name, p.avatar_url, p.is_admin, a.password_hash
		FROM auth_accounts a
		JOIN profiles p ON p.id = a.profile_id
		WHERE a.email = $1
	`, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.DisplayName,
		&user.AvatarURL,
		&user.IsAdmin,
		&passwordHash,
	)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)) != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := generateToken()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create session")
		return
	}
	expiresAt := time.Now().Add(sessionDuration)

	_, err = h.DB.Exec(r.Context(), `
		INSERT INTO sessions (profile_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`, user.ID, hashToken(token), expiresAt)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create session")
		return
	}

	setSessionCookie(w, token, expiresAt)
	writeJSON(w, http.StatusOK, user)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(sessionCookieName); err == nil && cookie.Value != "" {
		h.DB.Exec(r.Context(), `DELETE FROM sessions WHERE token_hash = $1`, hashToken(cookie.Value))
	}
	clearSessionCookie(w)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	user, ok := h.requireAuth(w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, user)
}
