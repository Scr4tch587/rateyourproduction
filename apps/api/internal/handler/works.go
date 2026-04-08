package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/scr4tch/rateyourproduction/apps/api/internal/model"
)

func (h *Handler) ListWorks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page := queryInt(r, "page", 1)
	perPage := queryInt(r, "per_page", 50)
	q := queryString(r, "q")
	if perPage > 100 {
		perPage = 100
	}
	offset := (page - 1) * perPage

	query := `
		SELECT w.id, w.slug, w.title, w.type, w.average_rating, w.rating_count, w.weighted_score,
			COALESCE(
				(SELECT array_agg(g.name) FROM work_genres wg JOIN genres g ON g.id = wg.genre_id WHERE wg.work_id = w.id),
				'{}'
			) AS genres,
			(SELECT count(*) FROM productions p WHERE p.work_id = w.id) AS production_count
		FROM works w
	`
	countQuery := `SELECT count(*) FROM works w`
	var args []any
	argIdx := 1

	if q != "" {
		query += ` WHERE w.title ILIKE $1 OR w.normalized_title ILIKE $1`
		countQuery += ` WHERE w.title ILIKE $1 OR w.normalized_title ILIKE $1`
		args = append(args, "%"+q+"%")
		argIdx++
	}

	query += ` ORDER BY w.weighted_score DESC`
	query += ` LIMIT $` + itoa(argIdx) + ` OFFSET $` + itoa(argIdx+1)
	args = append(args, perPage, offset)

	rows, err := h.DB.Query(ctx, query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query works")
		return
	}
	defer rows.Close()

	var works []model.WorkCard
	for rows.Next() {
		var wc model.WorkCard
		if err := rows.Scan(&wc.ID, &wc.Slug, &wc.Title, &wc.Type, &wc.AverageRating, &wc.RatingCount, &wc.WeightedScore, &wc.Genres, &wc.ProductionCount); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to scan work")
			return
		}
		works = append(works, wc)
	}

	var total int
	countArgs := args[:len(args)-2]
	if len(countArgs) > 0 {
		h.DB.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	} else {
		h.DB.QueryRow(ctx, countQuery).Scan(&total)
	}

	writeJSON(w, http.StatusOK, model.PaginatedResponse[model.WorkCard]{
		Data:    works,
		Total:   total,
		Page:    page,
		PerPage: perPage,
	})
}

func (h *Handler) GetWork(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slug := chi.URLParam(r, "slug")

	var work model.WorkDetail
	err := h.DB.QueryRow(ctx, `
		SELECT id, slug, title, normalized_title, type, description, premiere_year,
			average_rating, rating_count, weighted_score, created_at, updated_at
		FROM works WHERE slug = $1
	`, slug).Scan(
		&work.ID, &work.Slug, &work.Title, &work.NormalizedTitle, &work.Type,
		&work.Description, &work.PremiereYear, &work.AverageRating, &work.RatingCount,
		&work.WeightedScore, &work.CreatedAt, &work.UpdatedAt,
	)
	if err != nil {
		writeError(w, http.StatusNotFound, "work not found")
		return
	}

	// Genres
	genreRows, err := h.DB.Query(ctx, `
		SELECT g.id, g.name, g.slug
		FROM work_genres wg JOIN genres g ON g.id = wg.genre_id
		WHERE wg.work_id = $1
	`, work.ID)
	if err == nil {
		defer genreRows.Close()
		for genreRows.Next() {
			var g model.Genre
			genreRows.Scan(&g.ID, &g.Name, &g.Slug)
			work.Genres = append(work.Genres, g)
		}
	}

	// Creators
	creatorRows, err := h.DB.Query(ctx, `
		SELECT p.id, p.name, p.slug, wc.role_type
		FROM work_creators wc JOIN people p ON p.id = wc.person_id
		WHERE wc.work_id = $1
	`, work.ID)
	if err == nil {
		defer creatorRows.Close()
		for creatorRows.Next() {
			var c model.WorkCreator
			creatorRows.Scan(&c.PersonID, &c.Name, &c.Slug, &c.RoleType)
			work.Creators = append(work.Creators, c)
		}
	}

	// Productions
	prodRows, err := h.DB.Query(ctx, `
		SELECT p.id, p.work_id, p.slug, p.company_id, c.name, p.venue_id, v.name,
			p.city, p.country, p.start_date, p.end_date, p.production_label,
			p.average_rating, p.rating_count, p.weighted_score, p.created_at
		FROM productions p
		LEFT JOIN companies c ON c.id = p.company_id
		LEFT JOIN venues v ON v.id = p.venue_id
		WHERE p.work_id = $1
		ORDER BY p.start_date DESC NULLS LAST
	`, work.ID)
	if err == nil {
		defer prodRows.Close()
		for prodRows.Next() {
			var p model.Production
			prodRows.Scan(&p.ID, &p.WorkID, &p.Slug, &p.CompanyID, &p.CompanyName,
				&p.VenueID, &p.VenueName, &p.City, &p.Country, &p.StartDate, &p.EndDate,
				&p.ProductionLabel, &p.AverageRating, &p.RatingCount, &p.WeightedScore, &p.CreatedAt)
			work.Productions = append(work.Productions, p)
		}
	}

	writeJSON(w, http.StatusOK, work)
}
