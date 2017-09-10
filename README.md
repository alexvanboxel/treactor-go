




```
kubectl port-forward $(kubectl get pod -l app=grafana -o jsonpath='{.items[0].metadata.name}') 3000:3000 &
kubectl port-forward $(kubectl get pod -l app=servicegraph -o jsonpath='{.items[0].metadata.name}') 8088:8088 &
kubectl port-forward $(kubectl get pod -l app=zipkin -o jsonpath='{.items[0].metadata.name}') 9411:9411 &
kubectl port-forward $(kubectl get pod -l app=jaeger -o jsonpath='{.items[0].metadata.name}') 16686:16686 &

```



```bash
kubectl delete service reactor
kubectl delete service a
kubectl delete service b
kubectl delete service c

kubectl delete deploy reactor-v1
kubectl delete deploy reactor-service-a-v1
kubectl delete deploy reactor-service-b-v1
kubectl delete deploy reactor-service-c-v1
```