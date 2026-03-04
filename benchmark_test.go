package ull

import (
	"math/rand/v2"
	"testing"
)

func BenchmarkAdd(b *testing.B) {
	ull := MustNew(14)
	rng := rand.New(rand.NewPCG(42, 12345))

	// Pre-generate random values
	values := make([]uint64, b.N)
	for i := range values {
		values[i] = rng.Uint64()
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ull.Add(values[i])
	}
}

func BenchmarkAddBytes(b *testing.B) {
	ull := MustNew(14)

	data := []byte("benchmark test data for ultraloglog")

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		ull.AddBytes(data)
	}
}

func BenchmarkAddString(b *testing.B) {
	ull := MustNew(14)

	s := "benchmark test data for ultraloglog"

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		ull.AddString(s)
	}
}

func BenchmarkCount(b *testing.B) {
	benchmarks := []struct {
		name      string
		precision uint8
	}{
		{"precision_10", 10},
		{"precision_12", 12},
		{"precision_14", 14},
		{"precision_16", 16},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			ull := MustNew(bm.precision)
			rng := rand.New(rand.NewPCG(42, 12345))

			// Add some values
			for range 10000 {
				ull.Add(rng.Uint64())
			}

			b.ResetTimer()
			b.ReportAllocs()
			for b.Loop() {
				_ = ull.Count()
			}
		})
	}
}

func BenchmarkMerge(b *testing.B) {
	ull1 := MustNew(14)
	ull2 := MustNew(14)
	rng := rand.New(rand.NewPCG(42, 12345))

	for range 10000 {
		ull1.Add(rng.Uint64())
		ull2.Add(rng.Uint64())
	}

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		clone := ull1.Clone()
		_ = clone.Merge(ull2)
	}
}

func BenchmarkClone(b *testing.B) {
	benchmarks := []struct {
		name      string
		precision uint8
	}{
		{"precision_10", 10},
		{"precision_12", 12},
		{"precision_14", 14},
		{"precision_16", 16},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			ull := MustNew(bm.precision)
			rng := rand.New(rand.NewPCG(42, 12345))

			for range 10000 {
				ull.Add(rng.Uint64())
			}

			b.ResetTimer()
			b.ReportAllocs()
			for b.Loop() {
				_ = ull.Clone()
			}
		})
	}
}

func BenchmarkMarshalBinary(b *testing.B) {
	ull := MustNew(14)
	rng := rand.New(rand.NewPCG(42, 12345))

	for range 10000 {
		ull.Add(rng.Uint64())
	}

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_, _ = ull.MarshalBinary()
	}
}

func BenchmarkUnmarshalBinary(b *testing.B) {
	ull := MustNew(14)
	rng := rand.New(rand.NewPCG(42, 12345))

	for range 10000 {
		ull.Add(rng.Uint64())
	}

	data, _ := ull.MarshalBinary()

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		ull2 := &UltraLogLog{}
		_ = ull2.UnmarshalBinary(data)
	}
}

func BenchmarkHash64(b *testing.B) {
	benchmarks := []struct {
		name string
		size int
	}{
		{"8_bytes", 8},
		{"32_bytes", 32},
		{"128_bytes", 128},
		{"1024_bytes", 1024},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			data := make([]byte, bm.size)
			for i := range data {
				data[i] = byte(i)
			}

			b.ResetTimer()
			b.ReportAllocs()
			for b.Loop() {
				_ = hash64(data)
			}
		})
	}
}
