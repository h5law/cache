package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/h5law/cache/cache"
	"github.com/h5law/cache/workload"
)

func main() {
	capacity := flag.Int("capacity", 1024, "cache capacity")
	keys := flag.Uint64("keys", 100_000, "number of unique keys")
	operationCount := flag.Uint64("operations", 10_000_000, "number of operations")
	distribution := flag.String("distribution", "zipfian", "workload distribution")
	seed := flag.Uint64("seed", 42, "workload random seed")
	workingSetSize := flag.Uint64("working-set", 0, "working set size")

	flag.Parse()

	dist, err := parseDistribution(*distribution)
	if err != nil {
		panic(err)
	}

	generator := workload.NewGenerator(workload.GeneratorConfig{
		Seed:           *seed,
		Distribution:   dist,
		KeyCount:       *keys,
		Operations:     *operationCount,
		WorkingSetSize: *workingSetSize,
	})

	operations := generator.Generate()

	c := cache.New[uint64, uint64](*capacity)

	// Populate the cache before measuring the workload.
	for key := uint64(0); key < min(*keys, uint64(*capacity)); key++ {
		c.Put(key, key)
	}

	var hits uint64
	var misses uint64

	start := time.Now()

	for _, operation := range operations {
		switch operation.Type {
		case workload.OpGet:
			if _, ok := c.Get(operation.Key); ok {
				hits++
			} else {
				misses++
			}
		case workload.OpPut:
			c.Put(operation.Key, operation.Key)
		case workload.OpDelete:
			c.Delete(operation.Key)
		}
	}

	elapsed := time.Since(start)

	printResults(
		*capacity,
		*keys,
		*operationCount,
		*distribution,
		*seed,
		hits,
		misses,
		elapsed,
	)
}

func parseDistribution(value string) (workload.Distribution, error) {
	switch value {
	case "sequential":
		return workload.Sequential, nil
	case "uniform":
		return workload.Uniform, nil
	case "zipfian":
		return workload.Zipfian, nil
	case "hotcold":
		return workload.HotCold, nil
	case "scan":
		return workload.Scan, nil
	case "workingset":
		return workload.WorkingSet, nil
	default:
		return 0, fmt.Errorf("unknown distribution %q", value)
	}
}

func printResults(
	capacity int,
	keys uint64,
	operations uint64,
	distribution string,
	seed uint64,
	hits uint64,
	misses uint64,
	elapsed time.Duration,
) {
	requests := hits + misses

	var hitRate float64
	if requests > 0 {
		hitRate = float64(hits) / float64(requests) * 100
	}

	throughput := float64(requests) / elapsed.Seconds()

	fmt.Println("Cache Benchmark")
	fmt.Println("---------------")
	fmt.Printf("Policy:        LRU\n")
	fmt.Printf("Capacity:      %d\n", capacity)
	fmt.Printf("Keys:           %d\n", keys)
	fmt.Printf("Operations:     %d\n", operations)
	fmt.Printf("Distribution:   %s\n", distribution)
	fmt.Printf("Seed:           %d\n", seed)
	fmt.Println()
	fmt.Printf("Requests:       %d\n", requests)
	fmt.Printf("Hits:           %d\n", hits)
	fmt.Printf("Misses:         %d\n", misses)
	fmt.Printf("Hit rate:       %.2f%%\n", hitRate)
	fmt.Printf("Elapsed:        %s\n", elapsed)
	fmt.Printf("Throughput:     %.2f Mops/s\n", throughput/1_000_000)
}

func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}
