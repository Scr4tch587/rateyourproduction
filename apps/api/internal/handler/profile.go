package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/scr4tch/rateyourproduction/apps/api/internal/model"
)

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := chi.URLParam(r, "username")

	var profile model.Profile
	err := h.DB.QueryRow(ctx, `
		SELECT p.id, p.username, p.display_name, p.avatar_url, p.is_admin,
			(SELECT count(*) FROM logs l WHERE l.user_id = p.id) AS log_count,
			(SELECT count(*) FROM logs l WHERE l.user_id = p.id AND l.review_text IS NOT NULL) AS review_count
		FROM profiles p
		WHERE p.username = $1
	`, username).Scan(&profile.ID, &profile.Username, &profile.DisplayName, &profile.AvatarURL, &profile.IsAdmin, &profile.LogCount, &profile.ReviewCount)
	if err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	writeJSON(w, http.StatusOK, profile)
}
