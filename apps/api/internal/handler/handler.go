package handler

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/scr4tch/rateyourproduction/apps/api/internal/model"
)

type Handler struct {
	DB    *pgxpool.Pool
	Redis *redis.Client
}

func New(db *pgxpool.Pool, rdb *redis.Client) *Handler {
	return &Handler{DB: db, Redis: rdb}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func queryInt(r *http.Request, key string, fallback int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func queryFloat(r *http.Request, key string) *float64 {
	v := r.URL.Query().Get(key)
	if v == "" {
		return nil
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil
	}
	return &f
}

func queryIntPtr(r *http.Request, key string) *int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return nil
	}
	return &n
}

func querySlice(r *http.Request, key string) []string {
	values := r.URL.Query()[key]
	if len(values) == 0 {
		return nil
	}
	return values
}

func queryString(r *http.Request, key string) string {
	return strings.TrimSpace(r.URL.Query().Get(key))
}

const sessionCookieName = "ryp_session"

var errUnauthorized = errors.New("unauthorized")

func (h *Handler) currentUser(r *http.Request) (*model.SessionProfile, error) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil || cookie.Value == "" {
		return nil, errUnauthorized
	}

	tokenHash := hashToken(cookie.Value)

	var user model.SessionProfile
	err = h.DB.QueryRow(r.Context(), `
		SELECT p.id, p.username, a.email, p.display_name, p.avatar_url, p.is_admin
		FROM sessions s
		JOIN profiles p ON p.id = s.profile_id
		JOIN auth_accounts a ON a.profile_id = p.id
		WHERE s.token_hash = $1 AND s.expires_at > now()
	`, tokenHash).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.DisplayName,
		&user.AvatarURL,
		&user.IsAdmin,
	)
	if err != nil {
		return nil, errUnauthorized
	}

	return &user, nil
}

func (h *Handler) requireAuth(w http.ResponseWriter, r *http.Request) (*model.SessionProfile, bool) {
	user, err := h.currentUser(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return nil, false
	}
	return user, true
}

func (h *Handler) requireAdmin(w http.ResponseWriter, r *http.Request) (*model.SessionProfile, bool) {
	user, ok := h.requireAuth(w, r)
	if !ok {
		return nil, false
	}
	if !user.IsAdmin {
		writeError(w, http.StatusForbidden, "admin access required")
		return nil, false
	}
	return user, true
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func generateToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func setSessionCookie(w http.ResponseWriter, token string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  expiresAt,
	})
}

func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
}

func generateSlug(title string, city *string) string {
	s := strings.ToLower(strings.TrimSpace(title))
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "'", "")
	if city != nil && *city != "" {
		s += "-" + strings.ToLower(strings.ReplaceAll(strings.TrimSpace(*city), " ", "-"))
	}
	return s
}
