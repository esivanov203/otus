# Проверка результатов helm шаблонов
Кластер уже запущен, работает и протестирован, поэтому
```bash
helm template calendar ./k8s-charts --output-dir ./k8s-charts/generated
```
результирующие файлы ./k8s-charts/generated == с манифестами ./k8s через git diff