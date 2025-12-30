# Работоспособность системы

evg226@MacBook-Air-Evgeny hw12_13_14_15_16_calendar % curl -X POST http://calendar.local/events \
    -H "Content-Type: application/json" \
    -d '{
    "title": "K8s Homework Test",
    "description": "Testing full pipeline",
    "userId": "student-otus",
    "dateStart": "2024-12-30T10:00:00+03:00",
    "dateEnd": "2024-12-30T11:00:00+03:00"
    }'
fbf94a50-319a-488d-a5ad-c74cc1132b6f

evg226@MacBook-Air-Evgeny hw12_13_14_15_16_calendar % grpcurl -plaintext -import-path ./proto -proto calendar.proto \
    -d '{
        "event": {
        "title": "gRPC Test Event",
        "description": "Created via gRPC for homework",
        "user_id": "grpc-user-123",
        "date_start": 1735552800,
        "date_end": 1735556400
        }
    }' \
    grpc.calendar.local:80 calendar.CalendarService/CreateEvent
{
    "event": {
    "id": "2a65e69f-990e-4d44-88fc-a25c617a2237",
    "title": "gRPC Test Event",
    "description": "Created via gRPC for homework",
    "userId": "grpc-user-123",
    "dateStart": "1735552800",
    "dateEnd": "1735556400"
    }
}

evg226@MacBook-Air-Evgeny hw12_13_14_15_16_calendar % kubectl logs deployment/calendar --tail=4
{"level":"info","ts":"2025-12-30T11:20:04.159+0300","caller":"http/middleware.go:17","msg":"HTTP request","method":"GET","path":"/","proto":"HTTP/1.1","status":200,"latency":0,"userAgent":"curl/8.7.1","ip":"10.244.0.7:57126"}
{"level":"info","ts":"2025-12-30T11:40:37.949+0300","caller":"app/app.go:43","msg":"success created","app":"create event","title":"K8s Homework Test"}
{"level":"info","ts":"2025-12-30T11:40:37.949+0300","caller":"http/middleware.go:17","msg":"HTTP request","status":201,"latency":8,"userAgent":"curl/8.7.1","ip":"10.244.0.7:52426","method":"POST","path":"/events","proto":"HTTP/1.1"}
{"level":"info","ts":"2025-12-30T11:45:21.518+0300","caller":"app/app.go:43","msg":"success created","app":"create event","title":"gRPC Test Event"}

evg226@MacBook-Air-Evgeny hw12_13_14_15_16_calendar % kubectl logs deployment/scheduler --tail=4
{"level":"info","ts":"2025-12-30T11:18:38.684+0300","caller":"scheduler/scheduler_runner.go:59","msg":"scheduler is running..."}
{"level":"info","ts":"2025-12-30T11:40:38.633+0300","caller":"scheduler/scheduler.go:53","msg":"published event: fbf94a50-319a-488d-a5ad-c74cc1132b6f"}
{"level":"info","ts":"2025-12-30T11:45:28.616+0300","caller":"scheduler/scheduler.go:53","msg":"published event: 2a65e69f-990e-4d44-88fc-a25c617a2237"}

evg226@MacBook-Air-Evgeny hw12_13_14_15_16_calendar % kubectl logs deployment/sender --tail=4
{"level":"info","ts":"2025-12-30T11:18:38.719+0300","caller":"sender/sender_runner.go:51","msg":"sender is running..."}
{"level":"info","ts":"2025-12-30T11:40:38.636+0300","caller":"sender/sender.go:40","msg":"message has been received","id":"fbf94a50-319a-488d-a5ad-c74cc1132b6f"}
{"level":"info","ts":"2025-12-30T11:45:28.617+0300","caller":"sender/sender.go:40","msg":"message has been received","id":"2a65e69f-990e-4d44-88fc-a25c617a2237"}

evg226@MacBook-Air-Evgeny hw12_13_14_15_16_calendar % 
