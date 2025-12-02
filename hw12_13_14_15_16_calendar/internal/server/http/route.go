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
	s.router.HandleFunc("/events", s.createHandler).Methods(http.MethodPost)
	s.router.HandleFunc("/events", s.updateHandler).Methods(http.MethodPut)
	s.router.HandleFunc("/events/{id}", s.deleteHandler).Methods(http.MethodDelete)
	s.router.HandleFunc("/events/{id}", s.getOneHandler).Methods(http.MethodGet)
	s.router.HandleFunc("/events", s.getPeriodHandler).Methods(http.MethodGet)
}
