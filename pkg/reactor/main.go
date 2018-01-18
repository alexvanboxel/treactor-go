package main

import (
	"github.com/gorilla/mux"
	"net/http"
	//"fmt"
	"encoding/json"
	"github.com/alexvanboxel/reactor"
	"os"
	"fmt"
	"github.com/alexvanboxel/reactor/chem"
)

var port = os.Getenv("PORT")
var version = os.Getenv("VERSION")
var name = os.Getenv("NAME")
var mode = os.Getenv("MODE")

type ReactorInfo struct {
	Version     string
	VersionUser string
	VersionRole string
	VersionA    string
	VersionB    string
	VersionC    string
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{
		//CheckRedirect: redirectPolicyFunc,
	}
	trace := reactor.GetTrace(r)
	reactor.CallService(client, "http://role:3341/reactor/role", trace)

	version := reactor.Service{
		Version: version,
	}

	js, _ := json.MarshalIndent(version, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func RoleHandler(w http.ResponseWriter, r *http.Request) {
	version := reactor.Service{
		Version: version,
	}

	js, _ := json.MarshalIndent(version, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func WorkHandler(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{
		//CheckRedirect: redirectPolicyFunc,
	}
	trace := reactor.GetTrace(r)
	reactor.CallService(client, "http://user:3340/reactor/user", trace)
	reactor.CallService(client, "http://role:3341/reactor/role", trace)

	version := reactor.Service{
		Version: version,
	}

	js, _ := json.MarshalIndent(version, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func ReactorHandler(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{
		//CheckRedirect: redirectPolicyFunc,
	}

	//fmt.Println("========")
	//for k, v := range r.Header {
	//	fmt.Printf("%s = %s", k, v)
	//}

	trace := reactor.GetTrace(r)
	vu := reactor.CallService(client, "http://user:3340/reactor/user", trace)
	vr := reactor.CallService(client, "http://role:3341/reactor/role", trace)
	va := reactor.CallService(client, "http://a:3331/reactor/a", trace)
	vb := reactor.CallService(client, "http://b:3332/reactor/b", trace)
	vc := reactor.CallService(client, "http://c:3333/reactor/c", trace)

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

func main() {

	elemts := chem.ReadElements()

	fmt.Println(elemts)

	fmt.Printf("Reactor (%s:%s) listening on port %s\n", name, version, port)
	fmt.Printf("Mode: %s\n", mode)

	r := mux.NewRouter()
	r.HandleFunc("/reactor", ReactorHandler)
	r.HandleFunc("/reactor/a", WorkHandler)
	r.HandleFunc("/reactor/b", WorkHandler)
	r.HandleFunc("/reactor/c", WorkHandler)
	r.HandleFunc("/reactor/user", UserHandler)
	r.HandleFunc("/reactor/role", RoleHandler)
	http.Handle("/", r)

	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}
