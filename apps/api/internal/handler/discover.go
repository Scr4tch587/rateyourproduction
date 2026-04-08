package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/scr4tch/rateyourproduction/apps/api/internal/model"
)

func (h *Handler) Discover(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := model.DiscoverParams{
		Q:              queryString(r, "q"),
		Type:           parseWorkTypes(querySlice(r, "type")),
		GenreSlugs:     querySlice(r, "genre"),
		CompanySlugs:   querySlice(r, "company"),
		VenueSlugs:     querySlice(r, "venue"),
		City:           querySlice(r, "city"),
		Country:        querySlice(r, "country"),
		YearFrom:       queryIntPtr(r, "year_from"),
		YearTo:         queryIntPtr(r, "year_to"),
		MinRating:      queryFloat(r, "min_rating"),
		MinRatingCount: queryIntPtr(r, "min_rating_count"),
		CreatorSlugs:   querySlice(r, "creator"),
		PersonSlugs:    querySlice(r, "person"),
		RoleType:       querySlice(r, "role_type"),
		Sort:           r.URL.Query().Get("sort"),
		Page:           queryInt(r, "page", 1),
		PerPage:        queryInt(r, "per_page", 50),
	}
	if params.PerPage > 100 {
		params.PerPage = 100
	}

	var conditions []string
	var args []any
	argIdx := 1

	joins := ""
	needsProdJoin := false
	needsCreditJoin := false

	// Tier 1: work-level filters
	if params.Q != "" {
		conditions = append(conditions, fmt.Sprintf("(w.title ILIKE $%d OR w.normalized_title ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+params.Q+"%")
		argIdx++
	}
	if len(params.Type) > 0 {
		placeholders := makePlaceholders(&argIdx, len(params.Type))
		conditions = append(conditions, "w.type IN ("+placeholders+")")
		for _, t := range params.Type {
			args = append(args, string(t))
		}
	}
	if len(params.GenreSlugs) > 0 {
		placeholders := makePlaceholders(&argIdx, len(params.GenreSlugs))
		conditions = append(conditions, "EXISTS (SELECT 1 FROM work_genres wg JOIN genres g ON g.id = wg.genre_id WHERE wg.work_id = w.id AND g.slug IN ("+placeholders+"))")
		for _, s := range params.GenreSlugs {
			args = append(args, s)
		}
	}
	if params.MinRating != nil {
		conditions = append(conditions, fmt.Sprintf("w.average_rating >= $%d", argIdx))
		args = append(args, *params.MinRating)
		argIdx++
	}
	if params.MinRatingCount != nil {
		conditions = append(conditions, fmt.Sprintf("w.rating_count >= $%d", argIdx))
		args = append(args, *params.MinRatingCount)
		argIdx++
	}

	// Production-level filters require a join
	if len(params.CompanySlugs) > 0 || len(params.VenueSlugs) > 0 || len(params.City) > 0 || len(params.Country) > 0 || params.YearFrom != nil || params.YearTo != nil {
		needsProdJoin = true
	}
	if len(params.PersonSlugs) > 0 || len(params.RoleType) > 0 {
		needsProdJoin = true
		needsCreditJoin = true
	}

	if needsProdJoin {
		joins += " JOIN productions prod ON prod.work_id = w.id"

		if len(params.CompanySlugs) > 0 {
			joins += " LEFT JOIN companies comp ON comp.id = prod.company_id"
			placeholders := makePlaceholders(&argIdx, len(params.CompanySlugs))
			conditions = append(conditions, "comp.slug IN ("+placeholders+")")
			for _, s := range params.CompanySlugs {
				args = append(args, s)
			}
		}
		if len(params.VenueSlugs) > 0 {
			joins += " LEFT JOIN venues ven ON ven.id = prod.venue_id"
			placeholders := makePlaceholders(&argIdx, len(params.VenueSlugs))
			conditions = append(conditions, "ven.slug IN ("+placeholders+")")
			for _, s := range params.VenueSlugs {
				args = append(args, s)
			}
		}
		if len(params.City) > 0 {
			placeholders := makePlaceholders(&argIdx, len(params.City))
			conditions = append(conditions, "prod.city IN ("+placeholders+")")
			for _, s := range params.City {
				args = append(args, s)
			}
		}
		if len(params.Country) > 0 {
			placeholders := makePlaceholders(&argIdx, len(params.Country))
			conditions = append(conditions, "prod.country IN ("+placeholders+")")
			for _, s := range params.Country {
				args = append(args, s)
			}
		}
		if params.YearFrom != nil {
			conditions = append(conditions, fmt.Sprintf("EXTRACT(YEAR FROM prod.start_date) >= $%d", argIdx))
			args = append(args, *params.YearFrom)
			argIdx++
		}
		if params.YearTo != nil {
			conditions = append(conditions, fmt.Sprintf("EXTRACT(YEAR FROM prod.start_date) <= $%d", argIdx))
			args = append(args, *params.YearTo)
			argIdx++
		}
	}

	// Tier 2: creator filters
	if len(params.CreatorSlugs) > 0 {
		placeholders := makePlaceholders(&argIdx, len(params.CreatorSlugs))
		conditions = append(conditions, "EXISTS (SELECT 1 FROM work_creators wc JOIN people p ON p.id = wc.person_id WHERE wc.work_id = w.id AND p.slug IN ("+placeholders+"))")
		for _, s := range params.CreatorSlugs {
			args = append(args, s)
		}
	}

	// Tier 3: production credit filters
	if needsCreditJoin {
		joins += " JOIN production_credits pc ON pc.production_id = prod.id JOIN people cp ON cp.id = pc.person_id"
		if len(params.PersonSlugs) > 0 {
			placeholders := makePlaceholders(&argIdx, len(params.PersonSlugs))
			conditions = append(conditions, "cp.slug IN ("+placeholders+")")
			for _, s := range params.PersonSlugs {
				args = append(args, s)
			}
		}
		if len(params.RoleType) > 0 {
			placeholders := makePlaceholders(&argIdx, len(params.RoleType))
			conditions = append(conditions, "pc.role_type IN ("+placeholders+")")
			for _, s := range params.RoleType {
				args = append(args, s)
			}
		}
	}

	where := ""
	if len(conditions) > 0 {
		where = " WHERE " + strings.Join(conditions, " AND ")
	}

	orderBy := " ORDER BY w.weighted_score DESC"
	switch params.Sort {
	case "popular":
		orderBy = " ORDER BY w.rating_count DESC"
	case "newest":
		orderBy = " ORDER BY COALESCE((SELECT max(p3.start_date) FROM productions p3 WHERE p3.work_id = w.id), w.created_at::date) DESC"
	}

	offset := (params.Page - 1) * params.PerPage

	query := fmt.Sprintf(`
		SELECT DISTINCT w.id, w.slug, w.title, w.type, w.average_rating, w.rating_count, w.weighted_score,
			COALESCE(
				(SELECT array_agg(g.name) FROM work_genres wg JOIN genres g ON g.id = wg.genre_id WHERE wg.work_id = w.id),
				'{}'
			) AS genres,
			(SELECT count(*) FROM productions p2 WHERE p2.work_id = w.id) AS production_count
		FROM works w%s%s%s
		LIMIT $%d OFFSET $%d
	`, joins, where, orderBy, argIdx, argIdx+1)

	args = append(args, params.PerPage, offset)

	rows, err := h.DB.Query(ctx, query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query discover")
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

	// Count query
	countQuery := fmt.Sprintf(`SELECT count(DISTINCT w.id) FROM works w%s%s`, joins, where)
	countArgs := args[:len(args)-2]
	var total int
	if len(countArgs) > 0 {
		h.DB.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	} else {
		h.DB.QueryRow(ctx, countQuery).Scan(&total)
	}

	writeJSON(w, http.StatusOK, model.PaginatedResponse[model.WorkCard]{
		Data:    works,
		Total:   total,
		Page:    params.Page,
		PerPage: params.PerPage,
	})
}

func parseWorkTypes(raw []string) []model.WorkType {
	var types []model.WorkType
	for _, s := range raw {
		switch model.WorkType(s) {
		case model.WorkTypePlay, model.WorkTypeMusical, model.WorkTypeOpera:
			types = append(types, model.WorkType(s))
		}
	}
	return types
}

func makePlaceholders(argIdx *int, count int) string {
	parts := make([]string, count)
	for i := range count {
		parts[i] = fmt.Sprintf("$%d", *argIdx)
		*argIdx++
	}
	return strings.Join(parts, ", ")
}
