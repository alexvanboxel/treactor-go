package pi

import "testing"

func TestPiParallel(t *testing.T) {
	t.Logf("%f", Parallel(1*1000*1000))
}

func TestPiSingle(t *testing.T) {
	t.Logf("%f", Single(1*1000*1000))
}
