package reactor

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
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
	return trace
}

func CallService(client *http.Client, service string, trace map[string]string) (*Service) {
	req, _ := http.NewRequest("GET", service, nil)
	for k, v := range trace {
		req.Header.Set(k, v)
	}
	ra, err := client.Do(req)
	//ra, err := http.Get("http://a:3331/a")
	if err != nil {
		// handle error
	}
	defer ra.Body.Close()
	body, err := ioutil.ReadAll(ra.Body)

	va := Service{}
	err = json.Unmarshal(body, &va)
	return &va
}
