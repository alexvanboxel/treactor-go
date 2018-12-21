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

	//url := r.URL
	//plan, err := execute.Parse(url.Query().Get("atom"))
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//plan.Execute(r.Context())
	//_ = plan

	client.Logger.Info(r.Context(), "element log")
}

func Serve() {
	atoms := chem.NewAtoms()

	base := "/reactor"
	fmt.Printf("Reactor (%s:%s) listening on port %s\n", config.Name, config.Version, config.Port)
	fmt.Printf("Mode: %s\n", config.Mode)

	r := http.NewServeMux()
	r.HandleFunc(fmt.Sprintf("%s/split", base), ReactorSplit)
	r.HandleFunc(fmt.Sprintf("%s/orbit", base), ReactorOrbit)
	r.HandleFunc(fmt.Sprintf("%s/orbit/1", base), ReactorOrbit)
	r.HandleFunc(fmt.Sprintf("%s/orbit/2", base), ReactorOrbit)
	r.HandleFunc(fmt.Sprintf("%s/orbit/3", base), ReactorOrbit)
	r.HandleFunc(fmt.Sprintf("%s/orbit/4", base), ReactorOrbit)
	r.HandleFunc(fmt.Sprintf("%s/orbit/5", base), ReactorOrbit)
	r.HandleFunc(fmt.Sprintf("%s/orbit/inf", base), ReactorOrbit)
	for sym := range atoms.Symbols {
		r.HandleFunc(fmt.Sprintf("%s/atom/%s", base, sym), ReactorAtom)
	}
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Port), &ochttp.Handler{
		Handler:     r,
		Propagation: &propagation.HTTPFormat{},
	}))

}
