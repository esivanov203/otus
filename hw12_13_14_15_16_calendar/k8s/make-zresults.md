# Кластер
```bash
kubectl get all
```

evg226@MacBook-Air-Evgeny hw12_13_14_15_16_calendar % kubectl get all
NAME                                     READY   STATUS      RESTARTS   AGE
pod/calendar-79d6cd4b59-t9kp6            1/1     Running     0          20m
pod/calendar-migrations-4sgxn            0/1     Completed   0          21m
pod/calendar-postgres-6ff6dc989b-z9tvg   1/1     Running     0          22m
pod/calendar-rabbitmq-7599f9d7b8-mkv7w   1/1     Running     0          22m
pod/scheduler-54f6946566-ngdqj           1/1     Running     0          19m
pod/sender-85dcd77948-sn9cv              1/1     Running     0          19m

NAME                        TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)              AGE
service/calendar            ClusterIP   10.96.105.135   <none>        80/TCP,50051/TCP     20m
service/calendar-postgres   ClusterIP   10.96.239.186   <none>        5432/TCP             22m
service/calendar-rabbitmq   ClusterIP   10.96.200.64    <none>        5672/TCP,15672/TCP   22m
service/kubernetes          ClusterIP   10.96.0.1       <none>        443/TCP              26m

NAME                                READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/calendar            1/1     1            1           20m
deployment.apps/calendar-postgres   1/1     1            1           22m
deployment.apps/calendar-rabbitmq   1/1     1            1           22m
deployment.apps/scheduler           1/1     1            1           19m
deployment.apps/sender              1/1     1            1           19m

NAME                                           DESIRED   CURRENT   READY   AGE
replicaset.apps/calendar-79d6cd4b59            1         1         1       20m
replicaset.apps/calendar-postgres-6ff6dc989b   1         1         1       22m
replicaset.apps/calendar-rabbitmq-7599f9d7b8   1         1         1       22m
replicaset.apps/scheduler-54f6946566           1         1         1       19m
replicaset.apps/sender-85dcd77948              1         1         1       19m

NAME                            STATUS     COMPLETIONS   DURATION   AGE
job.batch/calendar-migrations   Complete   1/1           3s         21m
evg226@MacBook-Air-Evgeny hw12_13_14_15_16_calendar % 