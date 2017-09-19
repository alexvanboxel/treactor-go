




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





container-builder-local --config=cb-role.yaml --dryrun=false .
container-builder-local --config=cb-user.yaml --dryrun=false .
container-builder-local --config=cb-a.yaml --dryrun=false .
container-builder-local --config=cb-b.yaml --dryrun=false .
container-builder-local --config=cb-c.yaml --dryrun=false .
container-builder-local --config=cb-reactor.yaml --dryrun=false .




kubectl delete deploy reactor-role-v2
kubectl delete deploy reactor-user-v2
kubectl delete deploy reactor-service-a-v2
kubectl delete deploy reactor-service-b-v2
kubectl delete deploy reactor-service-c-v2
kubectl delete deploy reactor-v2

kubectl apply -f <(istioctl kube-inject -f k8s-reactor-v2.yaml)




kubectl delete deploy reactor-role-v1
kubectl delete deploy reactor-user-v1
kubectl delete deploy reactor-service-a-v1
kubectl delete deploy reactor-service-b-v1
kubectl delete deploy reactor-service-c-v1
kubectl delete deploy reactor-v1

kubectl apply -f <(istioctl kube-inject -f k8s-reactor.yaml)






istioctl create -f rule-2-role-segment-test.yaml
istioctl delete route-rules rule-2-role-segment-test

istioctl replace -f rule-3-user-weight50.yaml
istioctl delete route-rules rule-3-user-weight50
istioctl replace -f rule-4-user-v2.yaml
istioctl delete route-rules rule-6-b-delay

watch -n 1 -p -x curl -H "X-Segment: test" -s http://130.211.106.21/

