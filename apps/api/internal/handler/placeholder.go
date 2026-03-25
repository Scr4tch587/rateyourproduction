package handler

import (
	"encoding/json"
	"net/http"
)

func placeholder(resource string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"resource": resource,
			"status":   "not implemented",
		})
	}
}

func WorksList(w http.ResponseWriter, r *http.Request)       { placeholder("works")(w, r) }
func WorkGet(w http.ResponseWriter, r *http.Request)         { placeholder("work")(w, r) }
func ProductionsList(w http.ResponseWriter, r *http.Request)  { placeholder("productions")(w, r) }
func ProductionGet(w http.ResponseWriter, r *http.Request)    { placeholder("production")(w, r) }
func LogsList(w http.ResponseWriter, r *http.Request)         { placeholder("logs")(w, r) }
func LogCreate(w http.ResponseWriter, r *http.Request)        { placeholder("log_create")(w, r) }
func Discover(w http.ResponseWriter, r *http.Request)         { placeholder("discover")(w, r) }
func SubmissionCreate(w http.ResponseWriter, r *http.Request) { placeholder("submission_create")(w, r) }
func SubmissionsList(w http.ResponseWriter, r *http.Request)  { placeholder("submissions")(w, r) }
func AdminWorks(w http.ResponseWriter, r *http.Request)       { placeholder("admin_works")(w, r) }
func AdminProductions(w http.ResponseWriter, r *http.Request) { placeholder("admin_productions")(w, r) }
func AdminSubmissions(w http.ResponseWriter, r *http.Request) { placeholder("admin_submissions")(w, r) }
