package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/alexvanboxel/reactor/pkg/execute"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"encoding/json"
	"fmt"
	"github.com/alexvanboxel/reactor/pkg/chem"
	"github.com/alexvanboxel/reactor/pkg/client"
	"github.com/alexvanboxel/reactor/pkg/reactor"
	"log"
	"net/http"
	"os"
)

var port = os.Getenv("PORT")
var version = os.Getenv("VERSION")
var name = os.Getenv("NAME")
var mode = os.Getenv("MODE")

var gcpPubSub, _ = client.NewPubSub()

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
	capture := &RequestCapture{
		UrlPath:  r.URL.RawPath,
		UrlQuery: r.URL.RawQuery,
		Header:   headersFromRequest(r),
	}

	js, _ := json.MarshalIndent(capture, "", "\t")
	fmt.Println(string(js))

	msg := pubsub.Message{
		Data: js,
	}
	gcpPubSub.PostCaptured.Publish(context.Background(), &msg)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

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
		Version:     version,
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

	//fmt.Println(formula)
}

func main() {
	reactor.Configure()
	client.GoogleCloudInit()
	reactor.Init()

	elemts := chem.ReadElements()

	base := "/api/reactor"

	fmt.Println(elemts)

	fmt.Printf("Reactor (%s:%s) listening on port %s\n", name, version, port)
	fmt.Printf("Mode: %s\n", mode)

	r := http.NewServeMux()
	r.HandleFunc("/reactor/split", ReactorSplit)
	r.HandleFunc("/reactor/plane", ReactorPlane)
	r.HandleFunc("/reactor/element", ReactorElement)
	r.HandleFunc(fmt.Sprintf("%s/request/capture", base), HttpRequestCaptureHandler)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), &ochttp.Handler{
		Handler:     r,
		Propagation: &propagation.HTTPFormat{},
	}))

}
