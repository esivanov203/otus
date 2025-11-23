package internalhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/model"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"time"

	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/logger"
)

func (s *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	_ = r
	_, err := w.Write([]byte("Welcome to the root!"))
	if err != nil {
		s.logger.Error(err.Error(), logger.Fields{"handler": "root"})
	}
}

func (s *Server) decodeBodyToEvent(w http.ResponseWriter, r *http.Request, handlerName string) model.Event {
	var event model.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Error("decode JSON", logger.Fields{"handler": handlerName, "message": err.Error()})
		_, errW := w.Write([]byte("Invalid JSON payload"))
		if errW != nil {
			s.logger.Error(errW.Error(), logger.Fields{"handler": handlerName, "action": "send response"})
		}
	}
	return event
}

func (s *Server) createHandler(w http.ResponseWriter, r *http.Request) {
	event := s.decodeBodyToEvent(w, r, "create")
	if (event == model.Event{}) {
		return
	}

	event.ID = uuid.NewString()
	msg := event.ID

	err := s.app.CreateEvent(r.Context(), event)
	if err != nil {
		var ve *model.ValidationError
		if errors.As(err, &ve) {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		msg = err.Error()
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	_, errW := w.Write([]byte(msg))
	if errW != nil {
		s.logger.Error(errW.Error(), logger.Fields{"handler": "create", "action": "send response"})
	}
}

func (s *Server) updateHandler(w http.ResponseWriter, r *http.Request) {
	event := s.decodeBodyToEvent(w, r, "update")
	if (event == model.Event{}) {
		return
	}

	msg := event.ID

	err := s.app.UpdateEvent(r.Context(), event)
	if err != nil {
		var ve *model.ValidationError
		if errors.As(err, &ve) {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else if errors.Is(err, storage.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		msg = err.Error()
	} else {
		w.WriteHeader(http.StatusAccepted)
	}

	_, errW := w.Write([]byte(msg))
	if errW != nil {
		s.logger.Error(errW.Error(), logger.Fields{"handler": "update", "action": "send response"})
	}
}

func (s *Server) deleteHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := s.app.DeleteEvent(r.Context(), id)
	if err != nil {
		var ve *model.ValidationError
		if errors.As(err, &ve) {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else if errors.Is(err, storage.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_, errW := w.Write([]byte(err.Error()))
		if errW != nil {
			s.logger.Error(errW.Error(), logger.Fields{"handler": "delete", "action": "send response"})
		}
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Server) getOneHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	event, err := s.app.GetEvent(r.Context(), id)
	if err != nil {
		var ve *model.ValidationError
		if errors.As(err, &ve) {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else if errors.Is(err, storage.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_, errW := w.Write([]byte(err.Error()))
		if errW != nil {
			s.logger.Error(errW.Error(), logger.Fields{"handler": "getOne", "action": "send response"})
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(event); err != nil {
			s.logger.Error(err.Error(), logger.Fields{"handler": "getOne", "action": "encode response"})
		}
	}
}

func (s *Server) getPeriodHandler(w http.ResponseWriter, r *http.Request) {
	period := mux.Vars(r)["period"]
	userId := mux.Vars(r)["user_id"]
	startDate, _ := time.Parse("2006-01-02", mux.Vars(r)["start_date"])

	var (
		events       []model.Event
		err          error
		errBadPeriod = errors.New("invalid period")
	)

	switch period {
	case "day":
		events, err = s.app.ListEventsForDay(r.Context(), userId, startDate)
	case "week":
		events, err = s.app.ListEventsForWeek(r.Context(), userId, startDate)
	case "month":
		events, err = s.app.ListEventsForMonth(r.Context(), userId, startDate)
	default:
		err = fmt.Errorf("invalid period: %w", errBadPeriod)
	}

	if err != nil {
		var ve *model.ValidationError
		if errors.As(err, &ve) {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else if errors.Is(err, errBadPeriod) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_, errW := w.Write([]byte(err.Error()))
		if errW != nil {
			s.logger.Error(errW.Error(), logger.Fields{"handler": "getPeriod", "action": "send response"})
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(events); err != nil {
			s.logger.Error(err.Error(), logger.Fields{"handler": "getPeriods", "action": "encode response"})
		}
	}
}
