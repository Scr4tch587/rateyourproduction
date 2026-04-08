package model

import "time"

type WorkType string

const (
	WorkTypePlay    WorkType = "play"
	WorkTypeMusical WorkType = "musical"
	WorkTypeOpera   WorkType = "opera"
)

type Work struct {
	ID              string    `json:"id"`
	Slug            string    `json:"slug"`
	Title           string    `json:"title"`
	NormalizedTitle string    `json:"-"`
	Type            WorkType  `json:"type"`
	Description     *string   `json:"description,omitempty"`
	PremiereYear    *int      `json:"premiere_year,omitempty"`
	AverageRating   float64   `json:"average_rating"`
	RatingCount     int       `json:"rating_count"`
	WeightedScore   float64   `json:"weighted_score"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type WorkDetail struct {
	Work
	Genres      []Genre       `json:"genres"`
	Creators    []WorkCreator `json:"creators"`
	Productions []Production  `json:"productions,omitempty"`
}

type WorkCreator struct {
	PersonID string `json:"person_id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	RoleType string `json:"role_type"`
}

type WorkCard struct {
	ID              string   `json:"id"`
	Slug            string   `json:"slug"`
	Title           string   `json:"title"`
	Type            WorkType `json:"type"`
	Genres          []string `json:"genres"`
	AverageRating   float64  `json:"average_rating"`
	RatingCount     int      `json:"rating_count"`
	WeightedScore   float64  `json:"weighted_score"`
	ProductionCount int      `json:"production_count"`
}
