package internalhttp

import (
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/logger"
	"net/http"
)

func (s *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	_ = r
	_, err := w.Write([]byte("Welcome to the root!"))
	if err != nil {
		s.logger.Error(err.Error(), logger.Fields{"handler": "root"})
	}
}
