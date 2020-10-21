package libs

import (
	"github.com/pyto86pri/mackerel-lambda-agent/cmd/metrics"
)

// Map ...
func Map(vss []*metrics.Values) map[string][]float64 {
	vsm := make(map[string][]float64)
	for _, vs := range vss {
		for k, v := range *vs {
			vsm[k] = append(vsm[k], v)
		}
	}
	return vsm
}

func reduce(vs []float64, f func(float64, float64) float64) (r float64) {
	for _, v := range vs {
		r = f(r, v)
	}
	return
}

// Reduce ...
func Reduce(vm map[string][]float64, f func(float64, float64) float64) *metrics.Values {
	vs := make(metrics.Values)
	for k, vss := range vm {
		vs[k] = reduce(vss, f)
	}
	return &vs
}

// MapReduce ...
func MapReduce(vss []*metrics.Values, f func(float64, float64) float64) *metrics.Values {
	return Reduce(Map(vss), f)
}
