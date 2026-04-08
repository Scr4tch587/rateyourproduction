package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/scr4tch/rateyourproduction/apps/api/internal/model"
)

func (h *Handler) ListProductions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page := queryInt(r, "page", 1)
	perPage := queryInt(r, "per_page", 50)
	if perPage > 100 {
		perPage = 100
	}
	offset := (page - 1) * perPage
	workID := r.URL.Query().Get("work_id")

	query := `
		SELECT p.id, p.work_id, p.slug, p.company_id, c.name, p.venue_id, v.name,
			p.city, p.country, p.start_date, p.end_date, p.production_label,
			p.average_rating, p.rating_count, p.weighted_score, p.created_at
		FROM productions p
		LEFT JOIN companies c ON c.id = p.company_id
		LEFT JOIN venues v ON v.id = p.venue_id
	`
	countQuery := `SELECT count(*) FROM productions`
	var args []any
	argIdx := 1

	if workID != "" {
		query += ` WHERE p.work_id = $1`
		countQuery += ` WHERE work_id = $1`
		args = append(args, workID)
		argIdx++
	}

	query += ` ORDER BY p.start_date DESC NULLS LAST`
	query += ` LIMIT $` + itoa(argIdx) + ` OFFSET $` + itoa(argIdx+1)
	args = append(args, perPage, offset)

	rows, err := h.DB.Query(ctx, query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query productions")
		return
	}
	defer rows.Close()

	var productions []model.Production
	for rows.Next() {
		var p model.Production
		if err := rows.Scan(&p.ID, &p.WorkID, &p.Slug, &p.CompanyID, &p.CompanyName,
			&p.VenueID, &p.VenueName, &p.City, &p.Country, &p.StartDate, &p.EndDate,
			&p.ProductionLabel, &p.AverageRating, &p.RatingCount, &p.WeightedScore, &p.CreatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to scan production")
			return
		}
		productions = append(productions, p)
	}

	var total int
	if workID != "" {
		h.DB.QueryRow(ctx, countQuery, workID).Scan(&total)
	} else {
		h.DB.QueryRow(ctx, countQuery).Scan(&total)
	}

	writeJSON(w, http.StatusOK, model.PaginatedResponse[model.Production]{
		Data:    productions,
		Total:   total,
		Page:    page,
		PerPage: perPage,
	})
}

func (h *Handler) GetProduction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	var pd model.ProductionDetail
	err := h.DB.QueryRow(ctx, `
		SELECT p.id, p.work_id, p.slug, p.company_id, c.name, p.venue_id, v.name,
			p.city, p.country, p.start_date, p.end_date, p.production_label,
			p.average_rating, p.rating_count, p.weighted_score, p.created_at,
			w.title, w.slug
		FROM productions p
		LEFT JOIN companies c ON c.id = p.company_id
		LEFT JOIN venues v ON v.id = p.venue_id
		JOIN works w ON w.id = p.work_id
		WHERE p.id = $1
	`, id).Scan(
		&pd.ID, &pd.WorkID, &pd.Slug, &pd.CompanyID, &pd.CompanyName,
		&pd.VenueID, &pd.VenueName, &pd.City, &pd.Country, &pd.StartDate, &pd.EndDate,
		&pd.ProductionLabel, &pd.AverageRating, &pd.RatingCount, &pd.WeightedScore,
		&pd.CreatedAt, &pd.WorkTitle, &pd.WorkSlug,
	)
	if err != nil {
		writeError(w, http.StatusNotFound, "production not found")
		return
	}

	creditRows, err := h.DB.Query(ctx, `
		SELECT p.id, p.name, p.slug, pc.role_type
		FROM production_credits pc JOIN people p ON p.id = pc.person_id
		WHERE pc.production_id = $1
		ORDER BY pc.role_type, p.name
	`, pd.ID)
	if err == nil {
		defer creditRows.Close()
		for creditRows.Next() {
			var c model.ProductionCredit
			creditRows.Scan(&c.PersonID, &c.Name, &c.Slug, &c.RoleType)
			pd.Credits = append(pd.Credits, c)
		}
	}

	writeJSON(w, http.StatusOK, pd)
}

func itoa(n int) string {
	return strconv.Itoa(n)
}
