package internalgrpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/model"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *GRPCServer) Welcome(_ context.Context, _ *emptypb.Empty) (*proto.WelcomeResponse, error) {
	return &proto.WelcomeResponse{
		Message: "Welcome to Calendar gRPC",
	}, nil
}

func (s *GRPCServer) CreateEvent(ctx context.Context, req *proto.CreateEventRequest) (*proto.EventResponse, error) {
	event := model.Event{
		ID:          req.Event.Id,
		Title:       req.Event.Title,
		Description: req.Event.Description,
		UserID:      req.Event.UserId,
		DateStart:   time.Unix(req.Event.DateStart, 0),
		DateEnd:     time.Unix(req.Event.DateEnd, 0),
	}

	if event.ID == "" {
		event.ID = uuid.NewString()
	}

	err := s.app.CreateEvent(ctx, event)

	var ve model.ValidationError
	switch {
	case err == nil:
		return &proto.EventResponse{
			Event: &proto.Event{
				Id:          event.ID,
				Title:       event.Title,
				Description: event.Description,
				UserId:      event.UserID,
				DateStart:   event.DateStart.Unix(),
				DateEnd:     event.DateEnd.Unix(),
			},
		}, nil
	case errors.As(err, &ve):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	default:
		return nil, status.Error(codes.Internal, err.Error())
	}
}

func (s *GRPCServer) UpdateEvent(ctx context.Context, req *proto.UpdateEventRequest) (*proto.EventResponse, error) {
	event := model.Event{
		ID:          req.Event.Id,
		Title:       req.Event.Title,
		Description: req.Event.Description,
		UserID:      req.Event.UserId,
		DateStart:   time.Unix(req.Event.DateStart, 0),
		DateEnd:     time.Unix(req.Event.DateEnd, 0),
	}

	err := s.app.UpdateEvent(ctx, event)

	var ve model.ValidationError
	switch {
	case err == nil:
		return &proto.EventResponse{
			Event: &proto.Event{
				Id:          event.ID,
				Title:       event.Title,
				Description: event.Description,
				UserId:      event.UserID,
				DateStart:   event.DateStart.Unix(),
				DateEnd:     event.DateEnd.Unix(),
			},
		}, nil
	case errors.As(err, &ve):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, storage.ErrNotFound):
		return nil, status.Error(codes.NotFound, err.Error())
	default:
		return nil, status.Error(codes.Internal, err.Error())
	}
}

func (s *GRPCServer) DeleteEvent(
	ctx context.Context,
	req *proto.DeleteEventRequest,
) (*proto.DeleteEventResponse, error) {
	err := s.app.DeleteEvent(ctx, req.Id)

	var ve model.ValidationError
	switch {
	case err == nil:
		return &proto.DeleteEventResponse{}, nil
	case errors.As(err, &ve):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, storage.ErrNotFound):
		return nil, status.Error(codes.NotFound, err.Error())
	default:
		return nil, status.Error(codes.Internal, err.Error())
	}
}

func (s *GRPCServer) ListEventsPeriod(ctx context.Context, req *proto.ListRequest) (*proto.ListResponse, error) {
	startDate := time.Unix(req.Date, 0)

	var (
		events       []model.Event
		err          error
		errBadPeriod = errors.New("invalid period")
	)

	switch req.Period {
	case "day":
		events, err = s.app.ListEventsForDay(ctx, req.UserId, startDate)
	case "week":
		events, err = s.app.ListEventsForWeek(ctx, req.UserId, startDate)
	case "month":
		events, err = s.app.ListEventsForMonth(ctx, req.UserId, startDate)
	default:
		err = fmt.Errorf("param: %w", errBadPeriod)
	}

	var ve model.ValidationError
	switch {
	case err == nil:
		resp := &proto.ListResponse{}
		for _, ev := range events {
			resp.Event = append(resp.Event, &proto.Event{
				Id:          ev.ID,
				Title:       ev.Title,
				Description: ev.Description,
				UserId:      ev.UserID,
				DateStart:   ev.DateStart.Unix(),
				DateEnd:     ev.DateEnd.Unix(),
			})
		}
		return resp, nil
	case errors.As(err, &ve):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, errBadPeriod):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	default:
		return nil, status.Error(codes.Internal, err.Error())
	}
}

func (s *GRPCServer) GetEvent(ctx context.Context, req *proto.EventRequest) (*proto.EventResponse, error) {
	event, err := s.app.GetEvent(ctx, req.Id)

	var ve model.ValidationError
	switch {
	case err == nil:
		return &proto.EventResponse{
			Event: &proto.Event{
				Id:          event.ID,
				Title:       event.Title,
				Description: event.Description,
				UserId:      event.UserID,
				DateStart:   event.DateStart.Unix(),
				DateEnd:     event.DateEnd.Unix(),
			},
		}, nil
	case errors.As(err, &ve):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, storage.ErrNotFound):
		return nil, status.Error(codes.NotFound, err.Error())
	default:
		return nil, status.Error(codes.Internal, err.Error())
	}
}
