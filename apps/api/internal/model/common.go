package model

type Genre struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Person struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Company struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Slug    string  `json:"slug"`
	City    *string `json:"city,omitempty"`
	Country *string `json:"country,omitempty"`
}

type Venue struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Slug    string  `json:"slug"`
	City    *string `json:"city,omitempty"`
	Country *string `json:"country,omitempty"`
}

type Profile struct {
	ID          string  `json:"id"`
	Username    string  `json:"username"`
	DisplayName *string `json:"display_name,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	IsAdmin     bool    `json:"is_admin"`
	LogCount    int     `json:"log_count"`
	ReviewCount int     `json:"review_count"`
}

type PaginatedResponse[T any] struct {
	Data    []T `json:"data"`
	Total   int `json:"total"`
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}
