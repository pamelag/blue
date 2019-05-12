package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pamelag/blue/article"
	"github.com/pamelag/blue/tag"
)

const envString = "ART_ENV"

// Server holds the dependencies for a HTTP server.
type Server struct {
	Article article.Service
	Tag     tag.Service
	router  chi.Router
}

// New returns a new HTTP server.
func New(a article.Service, t tag.Service) *Server {
	s := &Server{
		Article: a,
		Tag:     t,
	}

	r := chi.NewRouter()
	r.Use(secureHeaders)
	r.Use(middleware.Recoverer)

	r.Route("/articles", func(r chi.Router) {
		ah := articleHandler{s.Article}
		r.Mount("/", ah.router())
	})

	r.Route("/tags", func(r chi.Router) {
		th := tagHandler{s.Tag}
		r.Mount("/", th.router())
	})

	s.router = r

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func secureHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		w.Header().Add("X-Frame-Options", "DENY")
		w.Header().Add("X-Content-Type-Options", "nosniff")
		w.Header().Add("X-XSS-Protection", "1; mode=block")
		w.Header().Add("Content-Security-Policy", "frame-ancestors 'none'")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case article.ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	case article.ErrArticleNotFound:
		w.WriteHeader(http.StatusNotFound)
	case tag.ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	case tag.ErrTagNotFound:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
