package handler

import (
	"encoding/json"
	"net/http"

	"github.com/scr4tch/rateyourproduction/apps/api/internal/model"
)

func (h *Handler) ListSubmissions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page := queryInt(r, "page", 1)
	perPage := queryInt(r, "per_page", 50)
	if perPage > 100 {
		perPage = 100
	}
	offset := (page - 1) * perPage
	status := queryString(r, "status")

	query := `
		SELECT s.id, s.work_id, w.title, s.submitted_by, p.username,
			s.company_id, c.name, s.venue_id, v.name, s.city, s.country,
			s.start_date, s.end_date, s.production_label, s.status, s.notes, s.created_at
		FROM production_submissions s
		JOIN works w ON w.id = s.work_id
		JOIN profiles p ON p.id = s.submitted_by
		LEFT JOIN companies c ON c.id = s.company_id
		LEFT JOIN venues v ON v.id = s.venue_id
	`
	countQuery := `SELECT count(*) FROM production_submissions s`
	var args []any
	argIdx := 1

	if status != "" {
		query += ` WHERE s.status = $1`
		countQuery += ` WHERE s.status = $1`
		args = append(args, status)
		argIdx++
	}

	query += ` ORDER BY s.created_at DESC`
	query += ` LIMIT $` + itoa(argIdx) + ` OFFSET $` + itoa(argIdx+1)
	args = append(args, perPage, offset)

	rows, err := h.DB.Query(ctx, query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query submissions")
		return
	}
	defer rows.Close()

	var submissions []model.ProductionSubmission
	for rows.Next() {
		var s model.ProductionSubmission
		if err := rows.Scan(
			&s.ID, &s.WorkID, &s.WorkTitle, &s.SubmittedBy, &s.SubmittedByName,
			&s.CompanyID, &s.CompanyName, &s.VenueID, &s.VenueName, &s.City, &s.Country,
			&s.StartDate, &s.EndDate, &s.ProductionLabel, &s.Status, &s.Notes, &s.CreatedAt,
		); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to scan submission")
			return
		}
		submissions = append(submissions, s)
	}

	var total int
	countArgs := args[:len(args)-2]
	if len(countArgs) > 0 {
		h.DB.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	} else {
		h.DB.QueryRow(ctx, countQuery).Scan(&total)
	}

	writeJSON(w, http.StatusOK, model.PaginatedResponse[model.ProductionSubmission]{
		Data:    submissions,
		Total:   total,
		Page:    page,
		PerPage: perPage,
	})
}

func (h *Handler) CreateSubmission(w http.ResponseWriter, r *http.Request) {
	user, ok := h.requireAuth(w, r)
	if !ok {
		return
	}

	var req model.SubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.WorkID == "" {
		writeError(w, http.StatusBadRequest, "work_id is required")
		return
	}

	var submissionID string
	err := h.DB.QueryRow(r.Context(), `
		INSERT INTO production_submissions (
			work_id, submitted_by, company_id, venue_id, city, country, start_date, end_date, production_label
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`, req.WorkID, user.ID, req.CompanyID, req.VenueID, req.City, req.Country, req.StartDate, req.EndDate, req.ProductionLabel).Scan(&submissionID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create submission")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"id": submissionID})
}
