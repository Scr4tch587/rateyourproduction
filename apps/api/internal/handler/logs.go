package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/scr4tch/rateyourproduction/apps/api/internal/model"
)

func (h *Handler) ListLogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page := queryInt(r, "page", 1)
	perPage := queryInt(r, "per_page", 50)
	if perPage > 100 {
		perPage = 100
	}
	offset := (page - 1) * perPage

	workID := r.URL.Query().Get("work_id")
	userID := r.URL.Query().Get("user_id")
	productionID := r.URL.Query().Get("production_id")

	query := `
		SELECT l.id, l.user_id, l.work_id, l.production_id, l.seen_date, l.rating, l.review_text, l.created_at,
			w.title, w.slug,
			p.production_label, c.name,
			pr.username
		FROM logs l
		JOIN works w ON w.id = l.work_id
		LEFT JOIN productions p ON p.id = l.production_id
		LEFT JOIN companies c ON c.id = p.company_id
		JOIN profiles pr ON pr.id = l.user_id
		WHERE 1=1
	`
	countQuery := `SELECT count(*) FROM logs WHERE 1=1`
	var args []any
	argIdx := 1

	if workID != "" {
		query += ` AND l.work_id = $` + itoa(argIdx)
		countQuery += ` AND work_id = $` + itoa(argIdx)
		args = append(args, workID)
		argIdx++
	}
	if userID != "" {
		query += ` AND l.user_id = $` + itoa(argIdx)
		countQuery += ` AND user_id = $` + itoa(argIdx)
		args = append(args, userID)
		argIdx++
	}
	if productionID != "" {
		query += ` AND l.production_id = $` + itoa(argIdx)
		countQuery += ` AND production_id = $` + itoa(argIdx)
		args = append(args, productionID)
		argIdx++
	}

	query += ` ORDER BY l.created_at DESC`
	query += ` LIMIT $` + itoa(argIdx) + ` OFFSET $` + itoa(argIdx+1)
	args = append(args, perPage, offset)

	rows, err := h.DB.Query(ctx, query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query logs")
		return
	}
	defer rows.Close()

	var logs []model.LogEntry
	for rows.Next() {
		var le model.LogEntry
		if err := rows.Scan(&le.ID, &le.UserID, &le.WorkID, &le.ProductionID,
			&le.SeenDate, &le.Rating, &le.ReviewText, &le.CreatedAt,
			&le.WorkTitle, &le.WorkSlug, &le.ProductionLabel, &le.CompanyName,
			&le.Username); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to scan log")
			return
		}
		logs = append(logs, le)
	}

	var total int
	countArgs := args[:len(args)-2]
	if len(countArgs) > 0 {
		h.DB.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	} else {
		h.DB.QueryRow(ctx, countQuery).Scan(&total)
	}

	writeJSON(w, http.StatusOK, model.PaginatedResponse[model.LogEntry]{
		Data:    logs,
		Total:   total,
		Page:    page,
		PerPage: perPage,
	})
}

func (h *Handler) CreateLog(w http.ResponseWriter, r *http.Request) {
	user, ok := h.requireAuth(w, r)
	if !ok {
		return
	}

	var req model.CreateLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.WorkID == "" {
		writeError(w, http.StatusBadRequest, "work_id is required")
		return
	}

	if req.Rating != nil && (*req.Rating < 0.5 || *req.Rating > 5.0) {
		writeError(w, http.StatusBadRequest, "rating must be between 0.5 and 5.0")
		return
	}

	var logID string
	err := h.DB.QueryRow(r.Context(), `
		INSERT INTO logs (user_id, work_id, production_id, seen_date, rating, review_text)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, user.ID, req.WorkID, req.ProductionID, req.SeenDate, req.Rating, req.ReviewText).Scan(&logID)

	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create log")
		return
	}

	// Update work rating aggregates
	h.updateWorkRating(r.Context(), req.WorkID)
	if req.ProductionID != nil {
		h.updateProductionRating(r.Context(), *req.ProductionID)
	}

	writeJSON(w, http.StatusCreated, map[string]string{"id": logID})
}

func (h *Handler) updateWorkRating(ctx context.Context, workID string) {
	// Bayesian weighted score: score = (v/(v+m))*R + (m/(v+m))*C
	// m = 10 (constant), C = global average
	h.DB.Exec(ctx, `
		WITH stats AS (
			SELECT
				COALESCE(AVG(rating), 0) AS avg_rating,
				COUNT(rating) AS rating_count
			FROM logs WHERE work_id = $1 AND rating IS NOT NULL
		),
		global AS (
			SELECT COALESCE(AVG(rating), 3.0) AS c FROM logs WHERE rating IS NOT NULL
		)
		UPDATE works SET
			average_rating = (SELECT avg_rating FROM stats),
			rating_count = (SELECT rating_count FROM stats),
			weighted_score = (
				SELECT
					CASE WHEN s.rating_count = 0 THEN 0
					ELSE (s.rating_count::numeric / (s.rating_count + 10)) * s.avg_rating
						+ (10::numeric / (s.rating_count + 10)) * g.c
					END
				FROM stats s, global g
			),
			updated_at = now()
		WHERE id = $1
	`, workID)
}

func (h *Handler) updateProductionRating(ctx context.Context, productionID string) {
	h.DB.Exec(ctx, `
		WITH stats AS (
			SELECT
				COALESCE(AVG(rating), 0) AS avg_rating,
				COUNT(rating) AS rating_count
			FROM logs WHERE production_id = $1 AND rating IS NOT NULL
		),
		global AS (
			SELECT COALESCE(AVG(rating), 3.0) AS c FROM logs WHERE rating IS NOT NULL
		)
		UPDATE productions SET
			average_rating = (SELECT avg_rating FROM stats),
			rating_count = (SELECT rating_count FROM stats),
			weighted_score = (
				SELECT
					CASE WHEN s.rating_count = 0 THEN 0
					ELSE (s.rating_count::numeric / (s.rating_count + 10)) * s.avg_rating
						+ (10::numeric / (s.rating_count + 10)) * g.c
					END
				FROM stats s, global g
			)
		WHERE id = $1
	`, productionID)
}
