package internalhttp

import (
	"github.com/gorilla/handlers"
	"net/http"
)

func (s *Server) setRouterUses() {
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.Use(s.loggingMiddleware)
}

func (s *Server) setRouterHandlers() {
	s.router.HandleFunc("/", s.rootHandler).Methods(http.MethodGet)
}
