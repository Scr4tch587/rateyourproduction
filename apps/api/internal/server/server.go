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

	r.Get("/health", handler.Health)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/works", func(r chi.Router) {
			r.Get("/", handler.WorksList)
			r.Get("/{slug}", handler.WorkGet)
		})

		r.Route("/productions", func(r chi.Router) {
			r.Get("/", handler.ProductionsList)
			r.Get("/{id}", handler.ProductionGet)
		})

		r.Route("/logs", func(r chi.Router) {
			r.Get("/", handler.LogsList)
			r.Post("/", handler.LogCreate)
		})

		r.Get("/discover", handler.Discover)

		r.Route("/submissions", func(r chi.Router) {
			r.Get("/", handler.SubmissionsList)
			r.Post("/", handler.SubmissionCreate)
		})

		r.Route("/admin", func(r chi.Router) {
			r.Get("/works", handler.AdminWorks)
			r.Get("/productions", handler.AdminProductions)
			r.Get("/submissions", handler.AdminSubmissions)
		})
	})

	return r
}
