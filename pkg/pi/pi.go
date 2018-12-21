package pi

import "math"

func Parallel(n int) float64 {
	ch := make(chan float64)
	for k := 0; k <= n; k++ {
		go termChannel(ch, float64(k))
	}
	f := 0.0
	for k := 0; k <= n; k++ {
		f += <-ch
	}
	return f
}

func termChannel(ch chan float64, k float64) {
	ch <- 4 * math.Pow(-1, k) / (2*k + 1)
}

func termFunc(k float64) float64 {
	return 4 * math.Pow(-1, k) / (2*k + 1)
}

func Single(n int) float64 {
	ch := make(chan float64)
	for k := 0; k <= n; k++ {
		go termChannel(ch, float64(k))
	}
	f := 0.0
	for k := 0; k <= n; k++ {
		f += termFunc(float64(k))
	}
	return f
}
