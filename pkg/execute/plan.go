package execute

type Plan interface {
}

type Repeat struct {
	times int
	mode  string
	block Plan
}

type KeyValue struct {
	kv map[string]string
}
