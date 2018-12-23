package reactor

import (
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"fmt"
	"github.com/alexvanboxel/reactor/pkg/chem"
	"github.com/alexvanboxel/reactor/pkg/client"
	"github.com/alexvanboxel/reactor/pkg/config"
	"github.com/alexvanboxel/reactor/pkg/execute"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"log"
	"net/http"
)

func ReactorSplit(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "Reactor.Split")
	defer span.End()
	//_, span := trace.StartSpan(r.Context(), "split.Get")
	//defer span.End()
	//span.Annotate([]trace.Attribute{trace.StringAttribute("key", "value")}, "something happened")
	//span.AddAttributes(trace.StringAttribute("hello", "world"))

	url := r.URL
	plan, err := execute.Parse(url.Query().Get("molecule"))
	if err != nil {
		fmt.Println(err)
		return
	}
	plan.Execute(ctx)

	client.Logger.Warning(ctx, "Test log")
}

func ReactorOrbit(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	plan, err := execute.Parse(url.Query().Get("molecule"))
	if err != nil {
		fmt.Println(err)
		return
	}
	plan.Execute(r.Context())
	client.Logger.Error(r.Context(), r, "Full error?")
}

func ReactorAtom(w http.ResponseWriter, r *http.Request) {
	//_, span := trace.StartSpan(r.Context(), "element.Get")
	//defer span.End()

	url := r.URL
	symbol := url.Query().Get("symbol")
	atom := client.Atoms.Symbols[symbol]
	//plan, err := execute.Parse(url.Query().Get("atom"))
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//plan.Execute(r.Context())
	//_ = plan

	client.Logger.Info(r.Context(), "Atom %s (%s)", atom.Name, atom.Number)
}

func ReactorHealthz(w http.ResponseWriter, r *http.Request) {
}

func Serve() {
	atoms := chem.NewAtoms()

	fmt.Printf("Reactor (%s:%s) listening on port %s\n", config.Name, config.Version, config.Port)
	fmt.Printf("Mode: %s\n", config.Mode)

	r := http.NewServeMux()
	r.HandleFunc("/", ReactorHealthz)
	r.HandleFunc("/healthz", ReactorHealthz)
	r.HandleFunc(fmt.Sprintf("%s/split", config.Base), ReactorSplit)
	r.HandleFunc(fmt.Sprintf("%s/orbit", config.Base), ReactorOrbit)
	r.HandleFunc(fmt.Sprintf("%s/orbit/1", config.Base), ReactorOrbit)
	r.HandleFunc(fmt.Sprintf("%s/orbit/2", config.Base), ReactorOrbit)
	r.HandleFunc(fmt.Sprintf("%s/orbit/3", config.Base), ReactorOrbit)
	r.HandleFunc(fmt.Sprintf("%s/orbit/4", config.Base), ReactorOrbit)
	r.HandleFunc(fmt.Sprintf("%s/orbit/5", config.Base), ReactorOrbit)
	r.HandleFunc(fmt.Sprintf("%s/orbit/inf", config.Base), ReactorOrbit)
	for sym := range atoms.Symbols {
		r.HandleFunc(fmt.Sprintf("%s/atom/%s", config.Base, sym), ReactorAtom)
	}
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Port), &ochttp.Handler{
		Handler:     r,
		Propagation: &propagation.HTTPFormat{},
	}))

}
