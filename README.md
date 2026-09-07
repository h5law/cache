# cache

## Motivation

Traditional cache replacement policies such as LRU and LFU rely on fixed
heuristics to decide which entries should be evicted.

This project explores a different approach:

> **Can a cache learn which entries are worth keeping?**

The project combines a conventional in-memory cache with a learned replacement
policy. The cache is responsible for storing and managing entries, whilst the
learned algorithm observes access patterns and predicts which entries are most
likely to be useful in the future.

The aim is to explore the intersection of **machine learning and systems
programming**, whilst keeping the implementation small enough to understand and
experiment with from first principles.

## Cache Design

The cache is split into two primary components:

```
                 ┌──────────────┐
                 │    Client    │
                 └──────┬───────┘
                        │
                     Get / Set
                        │
                 ┌──────▼───────┐
                 │     Cache    │
                 │              │
                 │   Storage    │
                 │   Metadata   │
                 │   Eviction   │
                 └──────┬───────┘
                        │
                 Replacement Policy
                        │
              ┌─────────▼─────────┐
              │      Learned      │
              │      Policy       │
              └───────────────────┘
```

The **cache** manages entries, lookups, insertions, deletions, capacity, and
eviction.

The **replacement policy** decides which entry should be removed when the cache
is full.

This separation allows the same cache implementation to be tested against
different policies:

```
Cache
 ├── Random
 ├── FIFO
 ├── LRU
 ├── LFU
 └── Learned
```

The learned policy can therefore be evaluated against traditional approaches
without changing the underlying cache.

## Project Architecture

Currently the cache project is structured in the following manner:

```
cache/
├── LICENSE
├── README.md
├── cache
│   ├── cache.go
│   ├── cache_test.go
│   └── lru.go
├── cmd
│   └── bench
│       └── main.go
├── go.mod
├── metrics
│   └── metrics.go
└── workload
    ├── generator.go
    ├── generator_test.go
    └── operation.go
```

The project is built in **Go**, with the initial implementation focusing on the
core cache and replacement-policy interface before introducing the learned
component.

## Future Extensions

The cache may eventually be extended beyond replacement to explore whether other
parts of the caching process can also benefit from learned algorithms.

Potential areas include:

- **Admission**
  — deciding if an item should be admitted into the cache in the first place.
- **Retrieval**
  — deciding how cached data should be located or retrieved efficiently.

This could eventually result in a cache where the three fundamental decisions
are independently optimised:

```
                         ┌──────────────┐
                         │    Request   │
                         └──────┬───────┘
                                │
                                ▼
                    ┌──────────────────────┐
                    │      Retrieval       │
                    │       Algorithm      │
                    └──────────┬───────────┘
                               │
                    ┌──────────┴───────────┐
                    │                      │
                  HIT                    MISS
                    │                      │
                    ▼                      ▼
               ┌─────────┐      ┌──────────────────────┐
               │  Value  │      │      Admission       │
               │ Returned│      │       Algorithm      │
               └─────────┘      └──────────┬───────────┘
                                           │
                                  ┌────────┴────────┐
                                  │                 │
                                REJECT            ADMIT
                                  │                 │
                                  ▼                 ▼
                              ┌────────┐    ┌───────────────┐
                              │  End   │    │     Cache     │
                              └────────┘    │    Storage    │
                                            └───────┬───────┘
                                                    │
                                                    │ Cache Full
                                                    ▼
                                           ┌──────────────────┐
                                           │    Replacement   │
                                           │     Algorithm    │
                                           └────────┬─────────┘
                                                    │
                                                    ▼
                                              Evict Entry
```

These components are intentionally not part of the initial implementation. They
provide potential directions for experimentation once the core cache and learned
replacement policy are established.
