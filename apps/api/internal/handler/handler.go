package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/scr4tch/rateyourproduction/apps/api/internal/config"
	"github.com/scr4tch/rateyourproduction/apps/api/internal/model"
)

type Handler struct {
	DB             *pgxpool.Pool
	Redis          *redis.Client
	SupabaseURL    string
	SupabaseAPIKey string
	HTTPClient     *http.Client
}

func New(cfg *config.Config, db *pgxpool.Pool, rdb *redis.Client) *Handler {
	apiKey := cfg.SupabaseAnonKey
	if apiKey == "" {
		apiKey = cfg.SupabaseServiceRoleKey
	}

	return &Handler{
		DB:             db,
		Redis:          rdb,
		SupabaseURL:    strings.TrimRight(cfg.SupabaseURL, "/"),
		SupabaseAPIKey: apiKey,
		HTTPClient:     &http.Client{},
	}
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

var errUnauthorized = errors.New("unauthorized")

func (h *Handler) currentUser(r *http.Request) (*model.SessionProfile, error) {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader == "" || !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		return nil, errUnauthorized
	}
	accessToken := strings.TrimSpace(authHeader[7:])
	if accessToken == "" {
		return nil, errUnauthorized
	}

	authUser, err := h.fetchSupabaseUser(r, accessToken)
	if err != nil {
		return nil, errUnauthorized
	}

	if err := h.ensureProfile(r, authUser); err != nil {
		return nil, errUnauthorized
	}

	var user model.SessionProfile
	err = h.DB.QueryRow(r.Context(), `
		SELECT p.id, p.username, p.display_name, p.avatar_url, p.is_admin
		FROM profiles p
		WHERE p.id = $1
	`, authUser.ID).Scan(
		&user.ID,
		&user.Username,
		&user.DisplayName,
		&user.AvatarURL,
		&user.IsAdmin,
	)
	if err != nil {
		return nil, errUnauthorized
	}
	user.Email = authUser.Email

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

func generateSlug(title string, city *string) string {
	s := strings.ToLower(strings.TrimSpace(title))
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "'", "")
	if city != nil && *city != "" {
		s += "-" + strings.ToLower(strings.ReplaceAll(strings.TrimSpace(*city), " ", "-"))
	}
	return s
}

type supabaseUser struct {
	ID           string         `json:"id"`
	Email        string         `json:"email"`
	UserMetadata map[string]any `json:"user_metadata"`
}

func (h *Handler) fetchSupabaseUser(r *http.Request, accessToken string) (*supabaseUser, error) {
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, h.SupabaseURL+"/auth/v1/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("apikey", h.SupabaseAPIKey)

	res, err := h.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		io.Copy(io.Discard, res.Body)
		return nil, fmt.Errorf("supabase auth user lookup failed: %d", res.StatusCode)
	}

	var user supabaseUser
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (h *Handler) ensureProfile(r *http.Request, authUser *supabaseUser) error {
	username, _ := authUser.UserMetadata["username"].(string)
	displayName, _ := authUser.UserMetadata["display_name"].(string)
	avatarURL, _ := authUser.UserMetadata["avatar_url"].(string)

	if strings.TrimSpace(username) == "" {
		localPart := strings.Split(strings.ToLower(authUser.Email), "@")[0]
		username = localPart + "-" + authUser.ID[:8]
	}

	_, err := h.DB.Exec(r.Context(), `
		INSERT INTO profiles (id, username, display_name, avatar_url, is_admin)
		VALUES (
			$1,
			$2,
			NULLIF($3, ''),
			NULLIF($4, ''),
			NOT EXISTS (SELECT 1 FROM profiles)
		)
		ON CONFLICT (id) DO UPDATE
		SET
			username = profiles.username,
			display_name = COALESCE(profiles.display_name, EXCLUDED.display_name),
			avatar_url = COALESCE(profiles.avatar_url, EXCLUDED.avatar_url)
	`, authUser.ID, username, displayName, avatarURL)
	return err
}
