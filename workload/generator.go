package workload

import "math/rand/v2"

type Distribution uint8

const (
	Sequential Distribution = iota
	Uniform
	Zipfian
	HotCold
	Scan
	WorkingSet
)

type GeneratorConfig struct {
	Seed           uint64
	Distribution   Distribution
	KeyCount       uint64
	Operations     uint64
	WorkingSetSize uint64
}

type Generator struct {
	rng          *rand.Rand
	distribution Distribution
	keyCount     uint64
	operations   uint64

	zipf *rand.Zipf

	workingSetSize uint64
	workingSetBase uint64
	workingSetPos  uint64
}

func NewGenerator(cfg GeneratorConfig) *Generator {
	if cfg.KeyCount == 0 {
		panic("workload: key count must be positive")
	}

	if cfg.Operations == 0 {
		panic("workload: operation count must be positive")
	}

	if cfg.WorkingSetSize == 0 {
		cfg.WorkingSetSize = cfg.KeyCount / 10
		if cfg.WorkingSetSize == 0 {
			cfg.WorkingSetSize = 1
		}
	}

	if cfg.WorkingSetSize > cfg.KeyCount {
		cfg.WorkingSetSize = cfg.KeyCount
	}

	rng := rand.New(
		rand.NewPCG(cfg.Seed, cfg.Seed^0x9e3779b97f4a7c15),
	)

	g := &Generator{
		rng:            rng,
		distribution:   cfg.Distribution,
		keyCount:       cfg.KeyCount,
		operations:     cfg.Operations,
		workingSetSize: cfg.WorkingSetSize,
	}

	if cfg.KeyCount > 1 {
		g.zipf = rand.NewZipf(rng, 1.1, 1, cfg.KeyCount-1)
	}

	return g
}

// Generate produces a deterministic sequence of GET operations.
func (g *Generator) Generate() []Operation {
	operations := make([]Operation, g.operations)

	for i := range operations {
		operations[i] = Operation{
			Type: OpGet,
			Key:  g.nextKey(uint64(i)),
		}
	}

	return operations
}

func (g *Generator) nextKey(operation uint64) uint64 {
	switch g.distribution {
	case Sequential:
		return operation % g.keyCount

	case Uniform:
		return g.rng.Uint64N(g.keyCount)

	case Zipfian:
		if g.keyCount == 1 {
			return 0
		}

		return g.zipf.Uint64()

	case HotCold:
		return g.hotColdKey()

	case Scan:
		return operation % g.keyCount

	case WorkingSet:
		return g.workingSetKey()

	default:
		panic("workload: unknown distribution")
	}
}

// hotColdKey generates a workload where 10% of keys receive 90% of requests.
func (g *Generator) hotColdKey() uint64 {
	hotKeys := g.keyCount / 10
	if hotKeys == 0 {
		hotKeys = 1
	}

	if g.rng.Float64() < 0.9 {
		return g.rng.Uint64N(hotKeys)
	}

	return g.rng.Uint64N(g.keyCount)
}

// workingSetKey generates requests within an active working set.
// The working set moves to a new random position after each complete pass.
func (g *Generator) workingSetKey() uint64 {
	if g.workingSetPos >= g.workingSetSize {
		g.workingSetPos = 0

		maxBase := g.keyCount - g.workingSetSize

		if maxBase == 0 {
			g.workingSetBase = 0
		} else {
			g.workingSetBase = g.rng.Uint64N(maxBase + 1)
		}
	}

	key := g.workingSetBase + g.workingSetPos
	g.workingSetPos++

	return key
}
