package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/scr4tch/rateyourproduction/apps/api/internal/config"
	"github.com/scr4tch/rateyourproduction/apps/api/internal/handler"
	mw "github.com/scr4tch/rateyourproduction/apps/api/internal/middleware"
)

type Server struct {
	cfg *config.Config
	db  *pgxpool.Pool
	rdb *redis.Client
}

func New(cfg *config.Config, db *pgxpool.Pool, rdb *redis.Client) *Server {
	return &Server{cfg: cfg, db: db, rdb: rdb}
}

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()

	r.Use(mw.CORS)
	r.Use(mw.Logger)
	r.Use(mw.Recoverer)
	r.Use(mw.Timeout)

	h := handler.New(s.db, s.rdb)

	r.Get("/health", handler.Health)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/signup", h.Signup)
			r.Post("/login", h.Login)
			r.Post("/logout", h.Logout)
			r.Get("/me", h.Me)
		})

		r.Route("/works", func(r chi.Router) {
			r.Get("/", h.ListWorks)
			r.Get("/{slug}", h.GetWork)
		})

		r.Route("/productions", func(r chi.Router) {
			r.Get("/", h.ListProductions)
			r.Get("/{id}", h.GetProduction)
		})

		r.Route("/logs", func(r chi.Router) {
			r.Get("/", h.ListLogs)
			r.Post("/", h.CreateLog)
		})

		r.Get("/discover", h.Discover)

		r.Route("/submissions", func(r chi.Router) {
			r.Get("/", h.ListSubmissions)
			r.Post("/", h.CreateSubmission)
		})

		r.Get("/profile/{username}", h.GetProfile)

		r.Route("/admin", func(r chi.Router) {
			r.Get("/works", h.AdminListWorks)
			r.Post("/works", h.AdminCreateWork)
			r.Delete("/works/{id}", h.AdminDeleteWork)
			r.Get("/productions", h.AdminListProductions)
			r.Delete("/productions/{id}", h.AdminDeleteProduction)
			r.Get("/submissions", h.AdminListSubmissions)
			r.Post("/submissions/{id}/approve", h.AdminApproveSubmission)
			r.Post("/submissions/{id}/reject", h.AdminRejectSubmission)
		})
	})

	return r
}
