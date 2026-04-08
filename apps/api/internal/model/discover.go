package model

type DiscoverParams struct {
	Q string `json:"q"`

	// Tier 1
	Type           []WorkType `json:"type"`
	GenreSlugs     []string   `json:"genre"`
	CompanySlugs   []string   `json:"company"`
	VenueSlugs     []string   `json:"venue"`
	City           []string   `json:"city"`
	Country        []string   `json:"country"`
	YearFrom       *int       `json:"year_from"`
	YearTo         *int       `json:"year_to"`
	MinRating      *float64   `json:"min_rating"`
	MinRatingCount *int       `json:"min_rating_count"`

	// Tier 2
	CreatorSlugs []string `json:"creator"`

	// Tier 3
	PersonSlugs []string `json:"person"`
	RoleType    []string `json:"role_type"`

	// Sorting & pagination
	Sort    string `json:"sort"`
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
}

type SubmissionRequest struct {
	WorkID          string  `json:"work_id"`
	CompanyID       *string `json:"company_id,omitempty"`
	VenueID         *string `json:"venue_id,omitempty"`
	City            *string `json:"city,omitempty"`
	Country         *string `json:"country,omitempty"`
	StartDate       *string `json:"start_date,omitempty"`
	EndDate         *string `json:"end_date,omitempty"`
	ProductionLabel *string `json:"production_label,omitempty"`
}
