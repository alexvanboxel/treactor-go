package main

import (
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"encoding/json"
	"fmt"
	"github.com/alexvanboxel/reactor/pkg/chem"
	"github.com/alexvanboxel/reactor/pkg/client"
	"github.com/alexvanboxel/reactor/pkg/config"
	"github.com/alexvanboxel/reactor/pkg/execute"
	"github.com/alexvanboxel/reactor/pkg/reactor"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"log"
	"net/http"
)

type RequestCapture struct {
	UrlPath  string            `json:"urlPath,string"`
	UrlQuery string            `json:"urlQuery,string"`
	Header   map[string]string `json:"header,string"`
	Body     string            `json:"url,string"`
}

type ReactorInfo struct {
	Version     string
	VersionUser string
	VersionRole string
	VersionA    string
	VersionB    string
	VersionC    string
}

func headersFromRequest(r *http.Request) (map[string]string) {
	calvinHeaders := make(map[string]string)
	for k, v := range r.Header {
		for i, vv := range v {
			calvinHeaders[fmt.Sprintf("%s.%d", k, i)] = vv
		}
	}
	return calvinHeaders
}

func HttpRequestCaptureHandler(w http.ResponseWriter, r *http.Request) {
	//capture := &RequestCapture{
	//	UrlPath:  r.URL.RawPath,
	//	UrlQuery: r.URL.RawQuery,
	//	Header:   headersFromRequest(r),
	//}
	//
	//js, _ := json.MarshalIndent(capture, "", "\t")
	//fmt.Println(string(js))
	//
	//msg := pubsub.Message{
	//	Data: js,
	//}
	////gcpPubSub.PostCaptured.Publish(context.Background(), &msg)
	//
	//w.Header().Set("Content-Type", "application/json")
	//w.Write(js)

}

func ReactorHandler(w http.ResponseWriter, r *http.Request) {

	//fmt.Println("========")
	//for k, v := range r.Header {
	//	fmt.Printf("%s = %s", k, v)
	//}

	trace := reactor.GetTrace(r)
	vu := reactor.CallService("http://user:3340/reactor/user", trace)
	vr := reactor.CallService("http://role:3341/reactor/role", trace)
	va := reactor.CallService("http://a:3331/reactor/a", trace)
	vb := reactor.CallService("http://b:3332/reactor/b", trace)
	vc := reactor.CallService("http://c:3333/reactor/c", trace)

	versions := ReactorInfo{
		Version:     config.Version,
		VersionUser: vu.Version,
		VersionRole: vr.Version,
		VersionA:    va.Version,
		VersionB:    vb.Version,
		VersionC:    vc.Version,
	}

	js, _ := json.MarshalIndent(versions, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func ReactorSplit(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "extra.Span")
	defer span.End()
	//_, span := trace.StartSpan(r.Context(), "split.Get")
	//defer span.End()
	//span.Annotate([]trace.Attribute{trace.StringAttribute("key", "value")}, "something happened")
	//span.AddAttributes(trace.StringAttribute("hello", "world"))

	url := r.URL
	_, _ = execute.Parse(url.Query().Get("molecule"))

	client.Logger.Warning(ctx, "Test log")
	reactor.CallPlane(ctx, 1, "2[U]")
	//fmt.Println(formula)
}

func ReactorPlane(w http.ResponseWriter, r *http.Request) {

	url := r.URL
	_, _ = execute.Parse(url.Query().Get("molecule"))

	reactor.CallElement(r.Context(), "U")
	//fmt.Println(formula)
}

func ReactorElement(w http.ResponseWriter, r *http.Request) {
	//_, span := trace.StartSpan(r.Context(), "element.Get")
	//defer span.End()

	url := r.URL
	_, _ = execute.Parse(url.Query().Get("molecule"))

	client.Logger.Info(r.Context(), "Element log")

	//fmt.Println(formula)
}

func main() {
	config.Configure()
	client.GoogleCloudInit()
	reactor.Init()

	elements := chem.ReadElements()
	fmt.Println(elements)

	base := "/reactor"
	fmt.Printf("Reactor (%s:%s) listening on port %s\n", config.Name, config.Version, config.Port)
	fmt.Printf("Mode: %s\n", config.Mode)

	r := http.NewServeMux()
	r.HandleFunc(fmt.Sprintf("%s/split", base), ReactorSplit)
	r.HandleFunc(fmt.Sprintf("%s/plane", base), ReactorPlane)
	r.HandleFunc(fmt.Sprintf("%s/element", base), ReactorElement)
	r.HandleFunc(fmt.Sprintf("%s/request/capture", base), HttpRequestCaptureHandler)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Port), &ochttp.Handler{
		Handler:     r,
		Propagation: &propagation.HTTPFormat{},
	}))

}
