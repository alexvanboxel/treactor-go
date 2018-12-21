package execute

import (
	"fmt"
	"github.com/alexvanboxel/reactor/pkg/client"
	"go.opencensus.io/trace"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

type Plan interface {
	String() string
	Execute(ctx context.Context)
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

func (o *Block) callElement(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	ctx, span := trace.StartSpan(ctx, "Reactor.Block.callElement")
	defer span.End()
	CallElement(ctx, o.block)
}

func (o *Block) callOrbit(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	ctx, span := trace.StartSpan(ctx, "Reactor.Block.callOrbit")
	defer span.End()
	CallOrbit(ctx, 1, o.block)
}

func (o *Block) Execute(ctx context.Context) {
	ctx, span := trace.StartSpan(ctx, "Reactor.Block")
	defer span.End()
	wg := sync.WaitGroup{}
	wg.Add(o.times)
	if o.mode == "s" {
		for i := 1; i <= o.times; i++ {
			if o.isAtom() {
				o.callElement(ctx, &wg)
			} else {
				o.callOrbit(ctx, &wg)
			}
		}
	} else if o.mode == "p" {
		for i := 1; i <= o.times; i++ {
			if o.isAtom() {
				go o.callElement(ctx, &wg)
			} else {
				go o.callOrbit(ctx, &wg)
			}
		}
	} else {
		// TODO ERR
	}
	wg.Wait()
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

func (o *Operator) Execute(ctx context.Context) {
	ctx, span := trace.StartSpan(ctx, "Reactor.Operator")
	defer span.End()
	wg := sync.WaitGroup{}
	wg.Add(2)
	if o.operand == PLUS {
		o.execute(ctx, &wg, o.left)
		o.execute(ctx, &wg, o.right)
	} else if o.operand == MULTIPLY {
		go o.execute(ctx, &wg, o.left)
		go o.execute(ctx, &wg, o.right)
	} else {
		// TODO ERR
	}
	wg.Wait()
}

func (o *Operator) execute(ctx context.Context, wg *sync.WaitGroup, plan Plan) {
	defer wg.Done()
	ctx, span := trace.StartSpan(ctx, "Reactor.Operator.execute")
	defer span.End()
	plan.Execute(ctx)
}

func (o *Operator) String() string {
	return o.left.String() + "?" + o.right.String()
}

func CallOrbit(context context.Context, orbit int, molecule string) {
	url := fmt.Sprintf("http://localhost:3330/reactor/orbit/%d?molecule=%s", orbit, molecule)
	req, _ := http.NewRequest("GET", url, nil)
	req = req.WithContext(context)
	ra, err := client.HttpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ra.Body.Close()
}

func CallElement(context context.Context, atom string) {
	full := atom
	atom = strings.Split(full, ",")[0]
	url := fmt.Sprintf("http://localhost:3330/reactor/atom/%s?atom=%s", atom, full)
	req, _ := http.NewRequest("GET", url, nil)
	req = req.WithContext(context)
	ra, err := client.HttpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ra.Body.Close()
}
