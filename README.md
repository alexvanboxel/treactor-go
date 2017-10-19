




```
kubectl -n istio-system port-forward $(kubectl -n istio-system get pod -l app=grafana -o jsonpath='{.items[0].metadata.name}') 3500:3000 &
kubectl port-forward -n istio-system $(kubectl get pod -n istio-system -l app=zipkin -o jsonpath='{.items[0].metadata.name}') 3501:9411 &
kubectl -n istio-system port-forward $(kubectl -n istio-system get pod -l app=servicegraph -o jsonpath='{.items[0].metadata.name}') 3502:8088 &   


kubectl -n istio-system port-forward $(kubectl -n istio-system get pod -l app=prometheus -o jsonpath='{.items[0].metadata.name}') 3509:9090 &

```





container-builder-local --config=cb.yaml --dryrun=false .


istioctl create -f rule-2-role-segment-test.yaml
istioctl delete route-rules rule-2-role-segment-test

istioctl replace -f rule-3-user-weight50.yaml
istioctl delete route-rules rule-3-user-weight50
istioctl replace -f rule-4-user-v2.yaml
istioctl delete route-rules rule-6-b-delay

watch -n 1 -p -x curl -s http://23.251.142.212/reactor
watch -n 1 -p -x curl -H "X-Segment: test" -s http://23.251.142.212/reactor

