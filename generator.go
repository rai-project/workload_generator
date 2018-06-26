package workload

import (
	"math/rand"
	"strings"

	"github.com/pkg/errors"
	"github.com/rai-project/micro18-tools/pkg/assets"

	"gonum.org/v1/gonum/stat/distuv"
)

// Pareto
// Zipf (not done)
// Uniform
// Exp
// Weibull
// Poisson

type Dist interface {
	// distuv.Quantiler
	distuv.RandLogProber
	CDF(float64) float64
	Prob(float64) float64
}

type Generator struct {
	done chan (bool)
	name string
	Dist
}

var (
	ValidDistributions = []string{
		"pareto",
		"zipf",
		"uniform",
		"exp", "exponential",
		"weibull",
		"poisson",
	}
	DefaultParetoParameters      = []float64{1, 1.5}
	DefaultZipfParameters        = []float64{}
	DefaultUniformParameters     = []float64{0, 1}
	DefaultExponentialParameters = []float64{0.5}
	DefaultWeibullParameters     = []float64{1.5, 1}
	DefaultPoissonParameters     = []float64{1}
)

func New(distribution string, params []float64) (*Generator, error) {
	if params == nil {
		params = []float64{}
	}
	var rnd *Generator
	switch strings.ToLower(distribution) {
	case "pareto":
		if len(params) != 2 {
			params = DefaultParetoParameters
		}
		rnd = NewPareto(params[0], params[1])
	case "zipf":
		return nil, errors.New("the zipf distribution is not implemented")
	case "uniform":
		if len(params) != 2 {
			params = DefaultUniformParameters
		}
		rnd = NewUniform(params[0], params[1])
	case "exp", "exponential":
		if len(params) != 1 {
			params = DefaultExponentialParameters
		}
		rnd = NewExponential(params[0])
	case "weibull":
		if len(params) != 2 {
			params = DefaultWeibullParameters
		}
		rnd = NewWeibull(params[0], params[1])
	case "poisson":
		if len(params) != 1 {
			params = DefaultPoissonParameters
		}
		rnd = NewPoisson(params[0])
	default:
		return nil, errors.Errorf("the distribution %s is unknown", distribution)
	}

	if rnd == nil {
		return nil, errors.Errorf("the distribution %s is unknown", distribution)
	}

	rnd.name = strings.ToLower(distribution)

	return rnd, nil
}

func NewUniform(min float64, max float64) *Generator {
	return &Generator{
		Dist: distuv.Uniform{Min: min, Max: max},
	}
}

func NewExponential(rate float64) *Generator {
	return &Generator{
		Dist: distuv.Exponential{Rate: rate},
	}
}

func NewWeibull(k float64, lambda float64) *Generator {
	return &Generator{
		Dist: distuv.Weibull{K: k, Lambda: lambda},
	}
}

func NewPoisson(lambda float64) *Generator {
	return &Generator{
		Dist: distuv.Poisson{Lambda: lambda},
	}
}

func NewZipf(...float64) *Generator {
	panic("not implemented")
	return nil
}

func NewPareto(xm float64, alpha float64) *Generator {
	return &Generator{
		Dist: distuv.Pareto{Xm: xm, Alpha: alpha},
	}
}

func (g *Generator) Next(at AliasTable, arry []interface{}) interface{} {
	idx := at.Next()
	return arry[idx]
}

func (g *Generator) Generator(arry []interface{}) <-chan interface{} {
	gen := make(chan interface{}, 10)
	at := NewAlias(g.probs(len(arry)), rand.NewSource(0))

	go func() {
		defer close(gen)
		for {
			select {
			case <-g.done:
				return
			default:
				gen <- g.Next(at, arry)
			}
		}
	}()
	return gen
}

func (g *Generator) probs(len int) []float64 {
	dist := g.Dist
	res := make([]float64, len)
	for ii := range res {
		res[ii] = dist.Rand() //float64(ii) * pmax / float64(len))
	}
	total := 0.0
	for _, r := range res {
		total = total + r
	}
	for ii := range res {
		res[ii] = res[ii] / total
	}
	return res
}

func (g *Generator) ModelGenerator(models assets.ModelManifests) <-chan assets.ModelManifest {
	gen := make(chan assets.ModelManifest, 10)
	arry := make([]interface{}, len(models))
	for ii, m := range models {
		arry[ii] = m
	}

	at := NewAlias(g.probs(len(arry)), rand.NewSource(0))

	go func() {
		defer close(gen)
		for {
			select {
			case <-g.done:
				return
			default:
				n := g.Next(at, arry)
				gen <- n.(assets.ModelManifest)
			}
		}
	}()
	return gen
}

func (g *Generator) Wait() {
	<-g.done
}

func (g *Generator) Close() {
	close(g.done)
}
