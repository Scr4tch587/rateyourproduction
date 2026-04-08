package model

import "time"

type Log struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	WorkID       string    `json:"work_id"`
	ProductionID *string   `json:"production_id,omitempty"`
	SeenDate     *string   `json:"seen_date,omitempty"`
	Rating       *float64  `json:"rating,omitempty"`
	ReviewText   *string   `json:"review_text,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type LogEntry struct {
	Log
	WorkTitle       string  `json:"work_title"`
	WorkSlug        string  `json:"work_slug"`
	ProductionLabel *string `json:"production_label,omitempty"`
	CompanyName     *string `json:"company_name,omitempty"`
	Username        string  `json:"username"`
}

type CreateLogRequest struct {
	WorkID       string   `json:"work_id"`
	ProductionID *string  `json:"production_id,omitempty"`
	SeenDate     *string  `json:"seen_date,omitempty"`
	Rating       *float64 `json:"rating,omitempty"`
	ReviewText   *string  `json:"review_text,omitempty"`
}
