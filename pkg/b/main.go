package main

import (
	"github.com/gorilla/mux"
	"net/http"
	//"fmt"
	"encoding/json"
	"github.com/alexvanboxel/reactor"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{
		//CheckRedirect: redirectPolicyFunc,
	}
	trace := reactor.GetTrace(r)
	reactor.CallService(client, "http://a:3340/reactor/user", trace)
	reactor.CallService(client, "http://a:3341/reactor/role", trace)

	version := reactor.Service{
		Version: "1",
	}

	js, _ := json.MarshalIndent(version, "","\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/reactor/b", HomeHandler)
	//r.HandleFunc("/products", ProductsHandler)
	//r.HandleFunc("/articles", ArticlesHandler)
	http.Handle("/", r)

	http.ListenAndServe(":3332", r)
}
