package workload

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkloadGenerator(t *testing.T) {
	// gen, err := New("uniform", []float64{})
	gen, err := New("pareto", []float64{})
	assert.NoError(t, err)
	len := 41
	lst := make([]interface{}, len)
	for ii := 0; ii < len; ii++ {
		lst[ii] = ii
	}
	at := NewAlias(gen.probs(len), rand.NewSource(0))
	for ii := 0; ii < 1000; ii++ {
		gen.Next(at, lst)
		// fmt.Println(r)
	}
}
