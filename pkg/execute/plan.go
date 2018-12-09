package execute

type Plan struct {
	Mode   string
	Blocks []Plan
}

