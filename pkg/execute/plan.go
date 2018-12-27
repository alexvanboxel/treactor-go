package execute

import (
	"encoding/json"
	"fmt"
	"github.com/alexvanboxel/reactor/pkg/resource"
	"github.com/alexvanboxel/reactor/pkg/config"
	"go.opencensus.io/trace"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

type Capture struct {
	Name     string
	Headers  map[string][]string
	Children []Capture
}

type Plan interface {
	String() string
	Execute(ctx context.Context, channel chan Capture)
	Calls() int
}

type Block struct {
	times int
	mode  string
	block string
	kv    map[string]string
}

func (o *Block) isAtom() bool {
	return unicode.IsLetter(rune(o.block[0]))
}

func (o *Block) callElement(ctx context.Context, wg *sync.WaitGroup, channel chan Capture) {
	defer wg.Done()
	// If REACTOR_TRACE_INTERNAL=1 add internal spans
	if config.TraceInternal() {
		context, span := trace.StartSpan(ctx, "Reactor.Block.callElement")
		defer span.End()
		ctx = context
	}
	CallElement(ctx, channel, o.block)
}

func (o *Block) callOrbit(ctx context.Context, wg *sync.WaitGroup, channel chan Capture) {
	defer wg.Done()
	// If REACTOR_TRACE_INTERNAL=1 add internal spans
	if config.TraceInternal() {
		context, span := trace.StartSpan(ctx, "Reactor.Block.callOrbit")
		defer span.End()
		ctx = context
	}
	CallOrbit(ctx, channel, o.block)
}

func (o *Block) Execute(ctx context.Context, channel chan Capture) {
	// If REACTOR_TRACE_INTERNAL=1 add internal spans
	if config.TraceInternal() {
		context, span := trace.StartSpan(ctx, "Reactor.Block")
		defer span.End()
		ctx = context
	}
	wg := sync.WaitGroup{}
	wg.Add(o.times)
	if o.mode == "s" {
		for i := 1; i <= o.times; i++ {
			if o.isAtom() {
				o.callElement(ctx, &wg, channel)
			} else {
				o.callOrbit(ctx, &wg, channel)
			}
		}
	} else if o.mode == "p" {
		for i := 1; i <= o.times; i++ {
			if o.isAtom() {
				go o.callElement(ctx, &wg, channel)
			} else {
				go o.callOrbit(ctx, &wg, channel)
			}
		}
	} else {
		// TODO ERR
	}
	wg.Wait()
}

func (o *Block) Calls() int  {
	return o.times
}

func (o *Block) String() string {
	s := strconv.Itoa(o.times) + o.mode + "[" + o.block + "]"
	for k, v := range o.kv {
		s += "," + k + ":" + v
	}
	return s
}

type Operator struct {
	left    Plan
	right   Plan
	operand Token
}

func (o *Operator) Execute(ctx context.Context, channel chan Capture) {
	// If REACTOR_TRACE_INTERNAL=1 add internal spans
	if config.TraceInternal() {
		context, span := trace.StartSpan(ctx, "Reactor.Operator")
		defer span.End()
		ctx = context
	}
	wg := sync.WaitGroup{}
	wg.Add(2)
	if o.operand == PLUS {
		o.execute(ctx, &wg, channel, o.left)
		o.execute(ctx, &wg, channel, o.right)
	} else if o.operand == MULTIPLY {
		go o.execute(ctx, &wg, channel, o.left)
		go o.execute(ctx, &wg, channel, o.right)
	} else {
		// TODO ERR
	}
	wg.Wait()
}

func (o *Operator) Calls() int  {
	return o.left.Calls() + o.right.Calls()
}

func (o *Operator) execute(ctx context.Context, wg *sync.WaitGroup, channel chan Capture, plan Plan) {
	defer wg.Done()
	// If REACTOR_TRACE_INTERNAL=1 add internal spans
	if config.TraceInternal() {
		context, span := trace.StartSpan(ctx, "Reactor.Operator.execute")
		defer span.End()
		ctx = context
	}
	plan.Execute(ctx, channel)
}

func (o *Operator) String() string {
	return o.left.String() + "?" + o.right.String()
}

func CallOrbit(context context.Context, channel chan Capture, molecule string) {
	next := config.NextOrbit()
	var url string
	if config.IsLocalMode() {
		url = fmt.Sprintf("http://localhost:%s%s/orbit/%s?molecule=%s", config.Port, config.Base, next, molecule)
	} else {
		url = fmt.Sprintf("http://orbit-%s%s/orbit/%s?molecule=%s", next, config.Base, next, molecule)
	}
	req, _ := http.NewRequest("GET", url, nil)
	req = req.WithContext(context)
	ra, err := resource.HttpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ra.Body.Close()
	decoder := json.NewDecoder(ra.Body)
	var capture Capture
	err = decoder.Decode(&capture)
	if err != nil {
		fmt.Println(err)
		return
	}
	channel <- capture
}

func CallElement(context context.Context, channel chan Capture, symbol string) {
	full := symbol
	symbol = strings.Split(full, ",")[0]
	var url string
	if config.IsLocalMode() {
		url = fmt.Sprintf("http://localhost:%s%s/atom/%s?symbol=%s", config.Port, config.Base, symbol, full)
	} else {
		url = fmt.Sprintf("http://atom-%s%s/atom/%s?symbol=%s", strings.ToLower(symbol), config.Base, symbol, full)
	}
	req, _ := http.NewRequest("GET", url, nil)
	req = req.WithContext(context)
	ra, err := resource.HttpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ra.Body.Close()
	decoder := json.NewDecoder(ra.Body)
	var capture Capture
	err = decoder.Decode(&capture)
	if err != nil {
		fmt.Println(err)
		return
	}
	channel <- capture
}
