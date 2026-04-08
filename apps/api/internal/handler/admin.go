package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/scr4tch/rateyourproduction/apps/api/internal/model"
)

type CreateWorkRequest struct {
	Title        string   `json:"title"`
	Type         string   `json:"type"`
	Description  *string  `json:"description,omitempty"`
	PremiereYear *int     `json:"premiere_year,omitempty"`
	GenreIDs     []string `json:"genre_ids,omitempty"`
	CreatorIDs   []struct {
		PersonID string `json:"person_id"`
		RoleType string `json:"role_type"`
	} `json:"creators,omitempty"`
}

func (h *Handler) AdminListWorks(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}
	h.ListWorks(w, r)
}

func (h *Handler) AdminCreateWork(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}
	var req CreateWorkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Title == "" || req.Type == "" {
		writeError(w, http.StatusBadRequest, "title and type are required")
		return
	}

	slug := strings.ToLower(strings.ReplaceAll(req.Title, " ", "-"))
	slug = strings.ReplaceAll(slug, "'", "")
	normalizedTitle := strings.ToLower(req.Title)

	ctx := r.Context()
	var workID string
	err := h.DB.QueryRow(ctx, `
		INSERT INTO works (slug, title, normalized_title, type, description, premiere_year)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, slug, req.Title, normalizedTitle, req.Type, req.Description, req.PremiereYear).Scan(&workID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create work")
		return
	}

	for _, gid := range req.GenreIDs {
		h.DB.Exec(ctx, `INSERT INTO work_genres (work_id, genre_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, workID, gid)
	}
	for _, c := range req.CreatorIDs {
		h.DB.Exec(ctx, `INSERT INTO work_creators (work_id, person_id, role_type) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`, workID, c.PersonID, c.RoleType)
	}

	writeJSON(w, http.StatusCreated, map[string]string{"id": workID})
}

func (h *Handler) AdminDeleteWork(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}
	id := chi.URLParam(r, "id")
	_, err := h.DB.Exec(r.Context(), `DELETE FROM works WHERE id = $1`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete work")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *Handler) AdminListProductions(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}
	h.ListProductions(w, r)
}

func (h *Handler) AdminDeleteProduction(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}
	id := chi.URLParam(r, "id")
	_, err := h.DB.Exec(r.Context(), `DELETE FROM productions WHERE id = $1`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete production")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *Handler) AdminListSubmissions(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.requireAdmin(w, r); !ok {
		return
	}
	h.ListSubmissions(w, r)
}

func (h *Handler) AdminApproveSubmission(w http.ResponseWriter, r *http.Request) {
	admin, ok := h.requireAdmin(w, r)
	if !ok {
		return
	}

	id := chi.URLParam(r, "id")
	var req model.ApproveSubmissionRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	ctx := r.Context()
	tx, err := h.DB.Begin(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to approve submission")
		return
	}
	defer tx.Rollback(ctx)

	var submission model.ProductionSubmission
	err = tx.QueryRow(ctx, `
		SELECT s.id, s.work_id, w.title, s.submitted_by, p.username,
			s.company_id, c.name, s.venue_id, v.name, s.city, s.country,
			s.start_date, s.end_date, s.production_label, s.status, s.notes, s.created_at
		FROM production_submissions s
		JOIN works w ON w.id = s.work_id
		JOIN profiles p ON p.id = s.submitted_by
		LEFT JOIN companies c ON c.id = s.company_id
		LEFT JOIN venues v ON v.id = s.venue_id
		WHERE s.id = $1
	`, id).Scan(
		&submission.ID, &submission.WorkID, &submission.WorkTitle, &submission.SubmittedBy, &submission.SubmittedByName,
		&submission.CompanyID, &submission.CompanyName, &submission.VenueID, &submission.VenueName, &submission.City, &submission.Country,
		&submission.StartDate, &submission.EndDate, &submission.ProductionLabel, &submission.Status, &submission.Notes, &submission.CreatedAt,
	)
	if err != nil {
		writeError(w, http.StatusNotFound, "submission not found")
		return
	}
	if submission.Status != model.SubmissionStatusPending {
		writeError(w, http.StatusBadRequest, "submission has already been reviewed")
		return
	}

	slug := generateSlug(submission.WorkTitle, submission.City)
	var productionID string
	err = tx.QueryRow(ctx, `
		INSERT INTO productions (work_id, slug, company_id, venue_id, city, country, start_date, end_date, production_label)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`, submission.WorkID, slug, submission.CompanyID, submission.VenueID, submission.City, submission.Country, submission.StartDate, submission.EndDate, submission.ProductionLabel).Scan(&productionID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create production")
		return
	}

	_, err = tx.Exec(ctx, `
		UPDATE production_submissions
		SET status = 'approved', notes = $2, reviewer_id = $3, reviewed_at = now()
		WHERE id = $1
	`, submission.ID, req.Notes, admin.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update submission")
		return
	}

	if err := tx.Commit(ctx); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to approve submission")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"id": productionID, "status": "approved"})
}

func (h *Handler) AdminRejectSubmission(w http.ResponseWriter, r *http.Request) {
	admin, ok := h.requireAdmin(w, r)
	if !ok {
		return
	}

	id := chi.URLParam(r, "id")
	var req model.ApproveSubmissionRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	tag, err := h.DB.Exec(r.Context(), `
		UPDATE production_submissions
		SET status = 'rejected', notes = $2, reviewer_id = $3, reviewed_at = now()
		WHERE id = $1 AND status = 'pending'
	`, id, req.Notes, admin.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to reject submission")
		return
	}
	if tag.RowsAffected() == 0 {
		writeError(w, http.StatusNotFound, "submission not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "rejected"})
}
