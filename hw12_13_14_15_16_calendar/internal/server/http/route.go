package internalhttp

import (
	"net/http"

	"github.com/gorilla/handlers"
)

func (s *Server) setRouterUses() {
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.Use(s.loggingMiddleware)
}

func (s *Server) setRouterHandlers() {
	s.router.HandleFunc("/", s.rootHandler).Methods(http.MethodGet)
}
