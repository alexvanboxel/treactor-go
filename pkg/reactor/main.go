package main

import (
	"github.com/gorilla/mux"
	"net/http"
	//"fmt"
	"encoding/json"
	"io/ioutil"
	"github.com/alexvanboxel/reactor"

)

type ReactorInfo struct {
	VersionA string
	VersionB string
	VersionC string
}


func HomeHandler(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "Category: %v\n", vars["category"])

	ra, err := http.Get("http://a:3331/a")
	if err != nil {
		// handle error
	}
	defer ra.Body.Close()
	body, err := ioutil.ReadAll(ra.Body)

	va := reactor.Service{}
	err = json.Unmarshal(body, &va)

	rb, err := http.Get("http://b:3332/b")
	if err != nil {
		// handle error
	}
	defer rb.Body.Close()
	body, err = ioutil.ReadAll(rb.Body)
	vb := reactor.Service{}
	err = json.Unmarshal(body, &vb)

	rc, err := http.Get("http://c:3333/c")
	if err != nil {
		// handle error
	}
	defer rc.Body.Close()
	body, err = ioutil.ReadAll(rc.Body)
	vc := reactor.Service{}
	err = json.Unmarshal(body, &vc)

	versions := ReactorInfo{
		VersionA:va.Version,
		VersionB:vb.Version,
		VersionC:vc.Version,
	}

	js, _ := json.MarshalIndent(versions, "","\t")
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
