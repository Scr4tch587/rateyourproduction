package model

import "time"

type SubmissionStatus string

const (
	SubmissionStatusPending  SubmissionStatus = "pending"
	SubmissionStatusApproved SubmissionStatus = "approved"
	SubmissionStatusRejected SubmissionStatus = "rejected"
)

type ProductionSubmission struct {
	ID              string           `json:"id"`
	WorkID          string           `json:"work_id"`
	WorkTitle       string           `json:"work_title"`
	SubmittedBy     string           `json:"submitted_by"`
	SubmittedByName string           `json:"submitted_by_name"`
	CompanyID       *string          `json:"company_id,omitempty"`
	CompanyName     *string          `json:"company_name,omitempty"`
	VenueID         *string          `json:"venue_id,omitempty"`
	VenueName       *string          `json:"venue_name,omitempty"`
	City            *string          `json:"city,omitempty"`
	Country         *string          `json:"country,omitempty"`
	StartDate       *string          `json:"start_date,omitempty"`
	EndDate         *string          `json:"end_date,omitempty"`
	ProductionLabel *string          `json:"production_label,omitempty"`
	Status          SubmissionStatus `json:"status"`
	Notes           *string          `json:"notes,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
}

type ApproveSubmissionRequest struct {
	Notes *string `json:"notes,omitempty"`
}
