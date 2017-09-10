package main

import (
	"github.com/gorilla/mux"
	"net/http"
	//"fmt"
	"encoding/json"
	"github.com/alexvanboxel/reactor"

	"fmt"
)

type ReactorInfo struct {
	VersionUser string
	VersionRole string
	VersionA    string
	VersionB    string
	VersionC    string
}


func HomeHandler(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{
		//CheckRedirect: redirectPolicyFunc,
	}

	fmt.Println("========")
	for k, v := range r.Header {
		fmt.Printf("%s = %s", k, v)
	}

	trace := reactor.GetTrace(r)
	vu := reactor.CallService(client, "http://user:3340/reactor/user", trace)
	vr := reactor.CallService(client, "http://role:3341/reactor/role", trace)
	va := reactor.CallService(client, "http://a:3331/reactor/a", trace)
	vb := reactor.CallService(client, "http://b:3332/reactor/b", trace)
	vc := reactor.CallService(client, "http://c:3333/reactor/c", trace)

	versions := ReactorInfo{
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

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	//r.HandleFunc("/products", ProductsHandler)
	//r.HandleFunc("/articles", ArticlesHandler)
	http.Handle("/", r)

	http.ListenAndServe(":3330", r)
}
