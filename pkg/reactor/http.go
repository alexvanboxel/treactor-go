package reactor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alexvanboxel/reactor/pkg/chem"
	"github.com/alexvanboxel/reactor/pkg/config"
	"github.com/alexvanboxel/reactor/pkg/execute"
	"github.com/alexvanboxel/reactor/pkg/resource"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"log"
	"net/http"
)

func executePlan(w http.ResponseWriter, r *http.Request, ctx context.Context, plan execute.Plan) {
	ch := make(chan execute.Capture, plan.Calls())
	plan.Execute(ctx, ch)

	elems := len(ch)
	capture := execute.Capture{
		Name:     config.AppName,
		Headers:  r.Header,
		Children: make([]execute.Capture, elems),
	}
	for i := 0; i < elems; i++ {
		capture.Children[i] = <-ch
	}

	bytes, _ := json.MarshalIndent(capture, "", "\t")
	w.Write(bytes)
}

func ReactorSplit(w http.ResponseWriter, r *http.Request) {
	// If REACTOR_TRACE_INTERNAL=1 add internal spans
	ctx := r.Context()
	if config.TraceInternal() {
		context, span := trace.StartSpan(r.Context(), "Reactor.Split")
		defer span.End()
		ctx = context
	}
	//_, span := trace.StartSpan(r.Context(), "split.Get")
	//defer span.End()
	//span.Annotate([]trace.Attribute{trace.StringAttribute("key", "value")}, "something happened")
	//span.AddAttributes(trace.StringAttribute("hello", "world"))
	url := r.URL
	molecule := url.Query().Get("molecule")
	resource.Logger.InfoF(ctx, "Starting reaction for molecule %s", molecule)

	plan, err := execute.Parse(molecule)
	if err != nil {
		resource.Logger.ErrorErr(ctx, r, "Unable to parse molecule", err)
		return
	}

	executePlan(w, r, ctx, plan)
	resource.Logger.WarningF(ctx, "Cooling down reaction, finished %s", molecule)
}

func ReactorOrbit(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	plan, err := execute.Parse(url.Query().Get("molecule"))
	if err != nil {
		fmt.Println(err)
		return
	}
	resource.Logger.Error(r.Context(), r, "Full error?")
	executePlan(w, r, r.Context(), plan)
}

func ReactorAtom(w http.ResponseWriter, r *http.Request) {
	//_, span := trace.StartSpan(r.Context(), "element.Get")
	//defer span.End()

	url := r.URL
	symbol := url.Query().Get("symbol")
	atom := resource.Atoms.Symbols[symbol]

	resource.Logger.InfoF(r.Context(), "Atom %s (%d)", atom.Name, atom.Number)
	capture := execute.Capture{
		Name:    config.AppName,
		Headers: r.Header,
	}
	bytes, _ := json.MarshalIndent(capture, "", "\t")
	w.Write(bytes)
}

func ReactorHealthz(w http.ResponseWriter, r *http.Request) {
}

func Serve() {
	atoms := chem.NewAtoms()

	fmt.Printf("Reactor (%s:%s) listening on port %s\n", config.AppName, config.AppVersion, config.Port)
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
		Propagation: config.TracePropagation(),
	}))
}
