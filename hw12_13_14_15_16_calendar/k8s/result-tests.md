evg226@MacBook-Air-Evgeny hw12_13_14_15_16_calendar % kubectl get all
NAME                                     READY   STATUS      RESTARTS   AGE
pod/calendar-79d6cd4b59-bzzmv            1/1     Running     0          60m
pod/calendar-migrations-vc7df            0/1     Completed   0          31m
pod/calendar-postgres-6ff6dc989b-lhfd5   1/1     Running     0          67m
pod/calendar-rabbitmq-7599f9d7b8-t49pf   1/1     Running     0          66m
pod/scheduler-54f6946566-ck62n           1/1     Running     0          29m
pod/sender-85dcd77948-ffm4s              1/1     Running     0          29m

NAME                        TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                        AGE
service/calendar            NodePort    10.96.60.27     <none>        80:30080/TCP,50051:30501/TCP   60m
service/calendar-postgres   ClusterIP   10.96.151.246   <none>        5432/TCP                       67m
service/calendar-rabbitmq   ClusterIP   10.96.139.140   <none>        5672/TCP,15672/TCP             66m
service/kubernetes          ClusterIP   10.96.0.1       <none>        443/TCP                        71m

NAME                                READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/calendar            1/1     1            1           60m
deployment.apps/calendar-postgres   1/1     1            1           67m
deployment.apps/calendar-rabbitmq   1/1     1            1           66m
deployment.apps/scheduler           1/1     1            1           29m
deployment.apps/sender              1/1     1            1           29m

NAME                                           DESIRED   CURRENT   READY   AGE
replicaset.apps/calendar-79d6cd4b59            1         1         1       60m
replicaset.apps/calendar-postgres-6ff6dc989b   1         1         1       67m
replicaset.apps/calendar-rabbitmq-7599f9d7b8   1         1         1       66m
replicaset.apps/scheduler-54f6946566           1         1         1       29m
replicaset.apps/sender-85dcd77948              1         1         1       29m

NAME                            STATUS     COMPLETIONS   DURATION   AGE
job.batch/calendar-migrations   Complete   1/1           3s         31m
evg226@MacBook-Air-Evgeny hw12_13_14_15_16_calendar % curl -X POST http://localhost/events \
-H "Content-Type: application/json" \
-d '{
"title": "Homework Test Event",
"description": "Created for K8s homework",
"userId": "student-123",
"dateStart": "2024-12-30T10:00:00+03:00",
"dateEnd": "2024-12-30T11:00:00+03:00"
}'
3577077e-ed1c-4ba5-b1dc-50fa84abd194%                                                                         evg226@MacBook-Air-Evgeny hw12_13_14_15_16_calendar % grpcurl -plaintext -import-path ./proto -proto calendar.proto \
-d '{
"event": {
"title": "Test Event",
"description": "Created via gRPC",
"user_id": "user123",
"date_start": 1735516800,
"date_end": 1735520400
}
}' \
localhost:50051 calendar.CalendarService/CreateEvent
{
"event": {
"id": "fea6d53d-65cf-44f1-97d3-d3a2ff0edd6a",
"title": "Test Event",
"description": "Created via gRPC",
"userId": "user123",
"dateStart": "1735516800",
"dateEnd": "1735520400"
}
}
evg226@MacBook-Air-Evgeny hw12_13_14_15_16_calendar % kubectl logs deployment/calendar --tail=2
{"level":"info","ts":"2025-12-29T17:50:08.614+0300","caller":"http/middleware.go:17","msg":"HTTP request","method":"POST","path":"/events","proto":"HTTP/1.1","status":201,"latency":2,"userAgent":"curl/8.7.1","ip":"10.244.0.1:1110"}
{"level":"info","ts":"2025-12-29T17:50:36.256+0300","caller":"app/app.go:43","msg":"success created","app":"create event","title":"Test Event"}
evg226@MacBook-Air-Evgeny hw12_13_14_15_16_calendar % kubectl logs deployment/calendar --tail=4
{"level":"info","ts":"2025-12-29T17:48:49.806+0300","caller":"http/middleware.go:17","msg":"HTTP request","status":201,"latency":1,"userAgent":"curl/8.7.1","ip":"10.244.0.1:52710","method":"POST","path":"/events","proto":"HTTP/1.1"}
{"level":"info","ts":"2025-12-29T17:50:08.614+0300","caller":"app/app.go:43","msg":"success created","app":"create event","title":"Homework Test Event"}
{"level":"info","ts":"2025-12-29T17:50:08.614+0300","caller":"http/middleware.go:17","msg":"HTTP request","method":"POST","path":"/events","proto":"HTTP/1.1","status":201,"latency":2,"userAgent":"curl/8.7.1","ip":"10.244.0.1:1110"}
{"level":"info","ts":"2025-12-29T17:50:36.256+0300","caller":"app/app.go:43","msg":"success created","app":"create event","title":"Test Event"}
evg226@MacBook-Air-Evgeny hw12_13_14_15_16_calendar % kubectl logs deployment/scheduler --tail=2
{"level":"info","ts":"2025-12-29T17:50:14.445+0300","caller":"scheduler/scheduler.go:53","msg":"published event: 3577077e-ed1c-4ba5-b1dc-50fa84abd194"}
{"level":"info","ts":"2025-12-29T17:50:44.437+0300","caller":"scheduler/scheduler.go:53","msg":"published event: fea6d53d-65cf-44f1-97d3-d3a2ff0edd6a"}
evg226@MacBook-Air-Evgeny hw12_13_14_15_16_calendar % kubectl logs deployment/sender --tail=2
{"level":"info","ts":"2025-12-29T17:50:14.446+0300","caller":"sender/sender.go:40","msg":"message has been received","id":"3577077e-ed1c-4ba5-b1dc-50fa84abd194"}
{"level":"info","ts":"2025-12-29T17:50:44.438+0300","caller":"sender/sender.go:40","msg":"message has been received","id":"fea6d53d-65cf-44f1-97d3-d3a2ff0edd6a"}
evg226@MacBook-Air-Evgeny hw12_13_14_15_16_calendar % 
