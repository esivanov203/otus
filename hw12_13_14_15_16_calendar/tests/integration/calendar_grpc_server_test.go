package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/esivanov203/otus/hw12_13_14_15_calendar/proto"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func newGRPCClient(t *testing.T) (proto.CalendarServiceClient, *grpc.ClientConn) {
	t.Helper()

	err := godotenv.Load("../../.env")
	require.NoError(t, err)
	host := os.Getenv("CALENDAR_HOST")
	port := "50051"
	addr := fmt.Sprintf("%s:%s", host, port)

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	client := proto.NewCalendarServiceClient(conn)
	return client, conn
}

func TestGrpcCrudEvent(t *testing.T) {
	client, conn := newGRPCClient(t)
	defer conn.Close()

	ctx := context.Background()

	// CREATE
	start := time.Now().Add(24 * time.Hour)
	end := time.Now().Add(7 * 24 * time.Hour)

	createResp, err := client.CreateEvent(ctx, &proto.CreateEventRequest{
		Event: &proto.Event{
			Title:       "New Event",
			Description: "New Event Description",
			UserId:      "u1-01",
			DateStart:   start.Unix(),
			DateEnd:     end.Unix(),
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, createResp.Event.Id)

	id := createResp.Event.Id

	// UPDATE
	_, err = client.UpdateEvent(ctx, &proto.UpdateEventRequest{
		Event: &proto.Event{
			Id:          id,
			Title:       "Updated Event",
			Description: "Updated",
			UserId:      "u1-01",
			DateStart:   start.Unix(),
			DateEnd:     end.Unix(),
		},
	})
	require.NoError(t, err)

	// GET ONE
	getResp, err := client.GetEvent(ctx, &proto.EventRequest{Id: id})
	require.NoError(t, err)
	require.Equal(t, "Updated Event", getResp.Event.Title)

	// LIST (week)
	listResp, err := client.ListEventsPeriod(ctx, &proto.ListRequest{
		UserId: "u1-01",
		Date:   time.Now().Unix(),
		Period: "week",
	})
	require.NoError(t, err)
	require.NotEmpty(t, listResp.Event)

	// DELETE
	_, err = client.DeleteEvent(ctx, &proto.DeleteEventRequest{Id: id})
	require.NoError(t, err)

	// GET AFTER DELETE → NOT_FOUND
	_, err = client.GetEvent(ctx, &proto.EventRequest{Id: id})
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, st.Code())
}

func TestGrpcCrudEventNegative(t *testing.T) {
	client, conn := newGRPCClient(t)
	defer conn.Close()

	ctx := context.Background()

	t.Run("create invalid dates", func(t *testing.T) {
		_, err := client.CreateEvent(ctx, &proto.CreateEventRequest{
			Event: &proto.Event{
				Title:     "Bad Event",
				UserId:    "u1-neg",
				DateStart: time.Now().Add(5 * time.Hour).Unix(),
				DateEnd:   time.Now().Add(1 * time.Hour).Unix(), // end < start
			},
		})

		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("update non-existing event", func(t *testing.T) {
		_, err := client.UpdateEvent(ctx, &proto.UpdateEventRequest{
			Event: &proto.Event{
				Id:        "550e8400-e29b-41d4-a716-446655440000",
				Title:     "Update",
				UserId:    "u1-neg",
				DateStart: time.Now().Unix(),
				DateEnd:   time.Now().Add(time.Hour).Unix(),
			},
		})

		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("get non-existing event", func(t *testing.T) {
		_, err := client.GetEvent(ctx, &proto.EventRequest{
			Id: "550e8400-e29b-41d4-a716-446655440000",
		})

		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("list invalid period", func(t *testing.T) {
		_, err := client.ListEventsPeriod(ctx, &proto.ListRequest{
			UserId: "u1-neg",
			Date:   time.Now().Unix(),
			Period: "WRONG_PERIOD",
		})

		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("list missing required fields", func(t *testing.T) {
		_, err := client.ListEventsPeriod(ctx, &proto.ListRequest{})
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.InvalidArgument, st.Code())
	})
}
