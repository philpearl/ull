# ull - UltraLogLog

[![GoDoc](https://godoc.org/github.com/philpearl/ull?status.svg)](https://godoc.org/github.com/philpearl/ull)

A Go implementation of UltraLogLog, a probabilistic cardinality estimation algorithm that improves upon HyperLogLog with better accuracy at the same memory cost.

## Installation

```bash
go get github.com/philpearl/ull
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/philpearl/ull"
)

func main() {
    // Create a new UltraLogLog with precision 14 (16KB memory, ~0.8% error)
    u := ull.MustNew(14)

    // Add elements
    u.AddString("hello")
    u.AddString("world")
    u.AddBytes([]byte("binary data"))

    // Or add pre-hashed values for better performance
    u.Add(0x123456789ABCDEF0)

    // Get the estimated cardinality
    count := u.Count()
    fmt.Printf("Estimated distinct elements: %d\n", count)
}
```

## Precision and Memory

| Precision | Memory | Standard Error |
|-----------|--------|----------------|
| 10        | 1 KB   | ~3.25%         |
| 12        | 4 KB   | ~1.625%        |
| 14        | 16 KB  | ~0.8125%       |
| 16        | 64 KB  | ~0.406%        |

## Features

- **High accuracy**: ~20% better accuracy than HyperLogLog
- **Configurable precision**: 4-18 bits (16 bytes to 256KB)
- **Mergeable**: Combine multiple UltraLogLogs for distributed counting
- **Serializable**: Binary marshaling for persistence
- **Fast**: Sub-nanosecond Add operations for pre-hashed values

## API

### Creating

```go
// With error handling
u, err := ull.New(14)

// Panic on invalid precision
u := ull.MustNew(14)
```

### Adding Elements

```go
u.Add(hash uint64)      // Add pre-hashed value (fastest)
u.AddBytes(data []byte) // Add bytes with internal hashing
u.AddString(s string)   // Add string with internal hashing
```

### Counting

```go
count := u.Count() // Returns estimated cardinality
```

### Merging

```go
u1 := ull.MustNew(14)
u2 := ull.MustNew(14)
// ... add elements to both ...

err := u1.Merge(u2) // u1 now contains the union
```

### Serialization

```go
// Marshal
data, err := u.MarshalBinary()

// Unmarshal
u2 := &ull.UltraLogLog{}
err = u2.UnmarshalBinary(data)
```

### Other Operations

```go
u.Clone()      // Create a deep copy
u.Reset()      // Clear all registers
u.Precision()  // Get the precision
u.Size()       // Get memory size in bytes
```

## Benchmarks

```
BenchmarkAdd             191587384     0.65 ns/op
BenchmarkAddBytes         17986622     6.69 ns/op
BenchmarkAddString        17998424     6.70 ns/op
BenchmarkCount/precision_14   1284    87469 ns/op
```

## License

MIT License - see [LICENSE](LICENSE) file.
