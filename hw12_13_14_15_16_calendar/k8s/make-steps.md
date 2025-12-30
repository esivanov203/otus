# Процесс создания кластера
```bash
kind create cluster --name calendar-cluster --config k8s/kind.yaml

kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml

kubectl create configmap calendar-env --from-env-file=.env
kubectl create configmap calendar-config --from-file=configs/config.yaml

kubectl apply -f k8s/postgres.yaml
kubectl apply -f k8s/rabbit.yaml

make build-img
kind load docker-image calendar-app:latest --name calendar-cluster
kind load docker-image calendar-scheduler:latest --name calendar-cluster
kind load docker-image calendar-sender:latest --name calendar-cluster
kind load docker-image calendar-migrations:latest --name calendar-cluster

kubectl apply -f k8s/migration.yaml
kubectl wait --for=condition=complete job/calendar-migrations --timeout=60s

kubectl apply -f k8s/sender.yaml
kubectl apply -f k8s/calendar.yaml
kubectl apply -f k8s/scheduler.yaml

kubectl apply -f k8s/ingress.yaml

echo "127.0.0.1 calendar.local" | sudo tee -a /etc/hosts
echo "127.0.0.1 grpc.calendar.local" | sudo tee -a /etc/hosts

kubectl get all
curl http://calendar.local/
grpcurl -plaintext -import-path ./proto -proto calendar.proto \
  -d '{}' grpc.calendar.local:80 calendar.CalendarService/Welcome
```