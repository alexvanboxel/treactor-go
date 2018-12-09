package reactor

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Service struct {
	Version string
}

func GetTrace(r *http.Request) (map[string]string) {
	trace := make(map[string]string, 0)
	trace["x-request-id"] = r.Header.Get("x-request-id")
	trace["x-b3-traceid"] = r.Header.Get("x-b3-traceid")
	trace["x-b3-spanid"] = r.Header.Get("x-b3-spanid")
	trace["x-b3-parentspanid"] = r.Header.Get("x-b3-parentspanid")
	trace["x-b3-sampled"] = r.Header.Get("x-b3-sampled")
	trace["x-b3-flags"] = r.Header.Get("x-b3-flags")
	trace["x-ot-span-context"] = r.Header.Get("x-ot-span-context")
	trace["x-segment"] = r.Header.Get("x-segment")
	return trace
}

func CallService(service string, trace map[string]string) (*Service) {
	req, _ := http.NewRequest("GET", service, nil)
	for k, v := range trace {
		req.Header.Set(k, v)
	}
	ra, err := HttpClient.Do(req)
	//ra, err := http.Get("http://a:3331/a")
	if err != nil {
		return &Service{
			Version: "-",
		}
		// handle error
	}
	defer ra.Body.Close()
	body, err := ioutil.ReadAll(ra.Body)

	va := Service{}
	err = json.Unmarshal(body, &va)
	return &va
}

func CallPlane(context context.Context, plane int, service string) {
	req, _ := http.NewRequest("GET", "http://localhost:3330/reactor/plane", nil)
	req = req.WithContext(context)
	ra, _ := HttpClient.Do(req)
	defer ra.Body.Close()
}

func CallElement(context context.Context, element string) {
	req, _ := http.NewRequest("GET", "http://localhost:3330/reactor/element", nil)
	req = req.WithContext(context)
	ra, _ := HttpClient.Do(req)
	defer ra.Body.Close()
}
