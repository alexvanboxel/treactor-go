# Reactor




### Molecule spec


```
S    [H,x=1,y=2,z=3]*2[O,x=1,y=2,z=3],x=1,y=2,z=3
A     H,x=1,y=2,z=3
A                      O,x=1,y=2,z=3
```


```
S    2[5[Ur,log:1,xyz:4]+5[C,log:1,xyz:4]],x:1,y:2
P1     5[Ur,log:1,xyz:4]+5[C,log:1,xyz:4]
A        Ur,log:1,xyz:4
A                          C,log:1,xyz:4
```


H,x=1,y=2,z=3


		//{"[Ur]", ""},
		//{"[Ur],log:1", ""},
		//{"5s[Ur]", ""},
		//{"5[Ur,log:1]", ""},
		//{"5[Ur,log:1,xyz:4]", ""},
		//{"5[Ur,log:1,xyz:4]+5[Ur,log:1,xyz:4]", ""},
		{"2[5[Ur,log:1,xyz:4]+5[Ur,log:1,xyz:4]],x:1,y:2", ""},






gcloud container clusters create istio \
    --enable-kubernetes-alpha \
    --machine-type=n1-standard-2 \
    --num-nodes=4 \
    --no-enable-legacy-authorization


curl -L https://git.io/getLatestIstio | sh -

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

watch -n 1 -p -x curl -s http://35.195.243.203/reactor
watch -n 1 -p -x curl -H "X-Segment: test" -s http://35.195.243.203/reactor

