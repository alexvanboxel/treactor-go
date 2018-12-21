package pi

import "testing"

func TestPiParallel(t *testing.T) {
	t.Logf("%f", Parallel(3*1000*1000))
}

func TestPiSingle(t *testing.T) {
	t.Logf("%f", Single(3*1000*1000))
}
