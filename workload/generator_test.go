package workload

import "testing"

func TestGeneratorDeterministic(t *testing.T) {
	config := GeneratorConfig{
		Seed:         42,
		Distribution: Uniform,
		KeyCount:     100,
		Operations:   1_000,
	}

	a := NewGenerator(config).Generate()
	b := NewGenerator(config).Generate()

	if len(a) != len(b) {
		t.Fatalf("expected equal lengths, got %d and %d", len(a), len(b))
	}

	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("workloads differ at operation %d: %+v != %+v", i, a[i], b[i])
		}
	}
}

func TestGeneratorProducesCorrectOperationCount(t *testing.T) {
	const operations = 1_000

	g := NewGenerator(GeneratorConfig{
		Seed:         42,
		Distribution: Uniform,
		KeyCount:     100,
		Operations:   operations,
	})

	workload := g.Generate()

	if len(workload) != operations {
		t.Fatalf("expected %d operations, got %d", operations, len(workload))
	}

	for i, operation := range workload {
		if operation.Type != OpGet {
			t.Fatalf("operation %d: expected OpGet, got %v", i, operation.Type)
		}
	}
}

func TestGeneratorKeysWithinBounds(t *testing.T) {
	distributions := []Distribution{
		Sequential,
		Uniform,
		Zipfian,
		HotCold,
		Scan,
		WorkingSet,
	}

	const (
		keyCount   = 100
		operations = 10_000
		workingSet = 20
	)

	for _, distribution := range distributions {
		t.Run(distributionName(distribution), func(t *testing.T) {
			g := NewGenerator(GeneratorConfig{
				Seed:           42,
				Distribution:   distribution,
				KeyCount:       keyCount,
				Operations:     operations,
				WorkingSetSize: workingSet,
			})

			for i, operation := range g.Generate() {
				if operation.Key >= keyCount {
					t.Fatalf(
						"operation %d: key %d outside [0, %d)",
						i,
						operation.Key,
						keyCount,
					)
				}
			}
		})
	}
}

func TestSequentialDistribution(t *testing.T) {
	g := NewGenerator(GeneratorConfig{
		Seed:         42,
		Distribution: Sequential,
		KeyCount:     5,
		Operations:   12,
	})

	workload := g.Generate()

	expected := []uint64{
		0, 1, 2, 3, 4,
		0, 1, 2, 3, 4,
		0, 1,
	}

	for i, expectedKey := range expected {
		if workload[i].Key != expectedKey {
			t.Fatalf(
				"operation %d: expected key %d, got %d",
				i,
				expectedKey,
				workload[i].Key,
			)
		}
	}
}

func TestScanDistribution(t *testing.T) {
	g := NewGenerator(GeneratorConfig{
		Seed:         42,
		Distribution: Scan,
		KeyCount:     4,
		Operations:   8,
	})

	workload := g.Generate()

	expected := []uint64{
		0, 1, 2, 3,
		0, 1, 2, 3,
	}

	for i, expectedKey := range expected {
		if workload[i].Key != expectedKey {
			t.Fatalf(
				"operation %d: expected key %d, got %d",
				i,
				expectedKey,
				workload[i].Key,
			)
		}
	}
}

func TestWorkingSetDistribution(t *testing.T) {
	g := NewGenerator(GeneratorConfig{
		Seed:           42,
		Distribution:   WorkingSet,
		KeyCount:       100,
		Operations:     1_000,
		WorkingSetSize: 10,
	})

	workload := g.Generate()

	for i := 0; i < len(workload); i += 10 {
		base := workload[i].Key

		for j := 0; j < 10 && i+j < len(workload); j++ {
			expected := base + uint64(j)

			if workload[i+j].Key != expected {
				t.Fatalf(
					"working set at operation %d: expected key %d, got %d",
					i+j,
					expected,
					workload[i+j].Key,
				)
			}
		}
	}
}

func TestSingleKeyWorkload(t *testing.T) {
	distributions := []Distribution{
		Sequential,
		Uniform,
		Zipfian,
		HotCold,
		Scan,
		WorkingSet,
	}

	for _, distribution := range distributions {
		t.Run(distributionName(distribution), func(t *testing.T) {
			g := NewGenerator(GeneratorConfig{
				Seed:         42,
				Distribution: distribution,
				KeyCount:     1,
				Operations:   100,
			})

			for i, operation := range g.Generate() {
				if operation.Key != 0 {
					t.Fatalf(
						"operation %d: expected key 0, got %d",
						i,
						operation.Key,
					)
				}
			}
		})
	}
}

func TestInvalidConfiguration(t *testing.T) {
	tests := []struct {
		name string
		cfg  GeneratorConfig
	}{
		{
			name: "zero keys",
			cfg: GeneratorConfig{
				KeyCount:   0,
				Operations: 100,
			},
		},
		{
			name: "zero operations",
			cfg: GeneratorConfig{
				KeyCount:   100,
				Operations: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("expected NewGenerator to panic")
				}
			}()

			NewGenerator(tt.cfg)
		})
	}
}

func distributionName(d Distribution) string {
	switch d {
	case Sequential:
		return "Sequential"
	case Uniform:
		return "Uniform"
	case Zipfian:
		return "Zipfian"
	case HotCold:
		return "HotCold"
	case Scan:
		return "Scan"
	case WorkingSet:
		return "WorkingSet"
	default:
		return "Unknown"
	}
}
