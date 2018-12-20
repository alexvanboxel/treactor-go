package execute

import "strconv"

type Plan interface {
	String() string
}

type Block struct {
	times int
	mode  string
	block string
	kv    map[string]string
}

type Operator struct {
	left    Plan
	right   Plan
	operand Token
}

func (o *Operator) String() string {
	return o.left.String() + "?" + o.right.String()
}

func (o *Block) String() string {
	s := strconv.Itoa(o.times) + o.mode + "[" + o.block + "]"
	for k, v := range o.kv {
		s += "," + k + ":" + v
	}
	return s
}
