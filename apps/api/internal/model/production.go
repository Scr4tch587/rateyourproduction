package model

import "time"

type Production struct {
	ID              string    `json:"id"`
	WorkID          string    `json:"work_id"`
	Slug            string    `json:"slug"`
	CompanyID       *string   `json:"company_id,omitempty"`
	CompanyName     *string   `json:"company_name,omitempty"`
	VenueID         *string   `json:"venue_id,omitempty"`
	VenueName       *string   `json:"venue_name,omitempty"`
	City            *string   `json:"city,omitempty"`
	Country         *string   `json:"country,omitempty"`
	StartDate       *string   `json:"start_date,omitempty"`
	EndDate         *string   `json:"end_date,omitempty"`
	ProductionLabel *string   `json:"production_label,omitempty"`
	AverageRating   float64   `json:"average_rating"`
	RatingCount     int       `json:"rating_count"`
	WeightedScore   float64   `json:"weighted_score"`
	CreatedAt       time.Time `json:"created_at"`
}

type ProductionDetail struct {
	Production
	WorkTitle string             `json:"work_title"`
	WorkSlug  string             `json:"work_slug"`
	Credits   []ProductionCredit `json:"credits"`
}

type ProductionCredit struct {
	PersonID string `json:"person_id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	RoleType string `json:"role_type"`
}
