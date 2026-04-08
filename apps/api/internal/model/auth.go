package model

type SignupRequest struct {
	Username    string  `json:"username"`
	Email       string  `json:"email"`
	Password    string  `json:"password"`
	DisplayName *string `json:"display_name,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SessionProfile struct {
	ID          string  `json:"id"`
	Username    string  `json:"username"`
	Email       string  `json:"email"`
	DisplayName *string `json:"display_name,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	IsAdmin     bool    `json:"is_admin"`
}
