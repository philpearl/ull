package ull

import (
	"fmt"
	"math"
	"math/rand/v2"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		precision uint8
		wantErr   bool
	}{
		{3, true},   // Too small
		{4, false},  // Minimum valid
		{12, false}, // Common value
		{18, false}, // Maximum valid
		{19, true},  // Too large
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("precision_%d", tt.precision), func(t *testing.T) {
			ull, err := New(tt.precision)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ull.Precision() != tt.precision {
				t.Errorf("precision = %d, want %d", ull.Precision(), tt.precision)
			}
			if ull.Size() != 1<<tt.precision {
				t.Errorf("size = %d, want %d", ull.Size(), 1<<tt.precision)
			}
		})
	}
}

func TestMustNew(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		ull := MustNew(12)
		if ull.Precision() != 12 {
			t.Errorf("precision = %d, want 12", ull.Precision())
		}
	})

	t.Run("panics_on_invalid", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic, got none")
			}
		}()
		MustNew(3)
	})
}

func TestAdd(t *testing.T) {
	ull := MustNew(14)

	// Add nothing, count should be 0
	if count := ull.Count(); count != 0 {
		t.Errorf("empty count = %d, want 0", count)
	}

	// Add one element
	ull.Add(0x123456789ABCDEF0)
	if count := ull.Count(); count == 0 {
		t.Error("count should be > 0 after adding element")
	}
}

func TestAddBytes(t *testing.T) {
	ull := MustNew(14)

	ull.AddBytes([]byte("hello"))
	ull.AddBytes([]byte("world"))

	count := ull.Count()
	if count < 1 || count > 5 {
		t.Errorf("count = %d, expected 1-5 for 2 distinct elements", count)
	}
}

func TestAddString(t *testing.T) {
	ull := MustNew(14)

	ull.AddString("hello")
	ull.AddString("world")

	count := ull.Count()
	if count < 1 || count > 5 {
		t.Errorf("count = %d, expected 1-5 for 2 distinct elements", count)
	}
}

func TestCardinalityEstimation(t *testing.T) {
	tests := []struct {
		name      string
		n         int
		precision uint8
		tolerance float64 // Acceptable error as a fraction
	}{
		{"small_100", 100, 14, 0.10},
		{"medium_1000", 1000, 14, 0.05},
		{"large_10000", 10000, 14, 0.03},
		{"very_large_100000", 100000, 14, 0.02},
		{"low_precision_1000", 1000, 10, 0.10},
		{"high_precision_1000", 1000, 16, 0.03},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ull := MustNew(tt.precision)
			rng := rand.New(rand.NewPCG(42, 12345))

			for i := 0; i < tt.n; i++ {
				ull.Add(rng.Uint64())
			}

			estimate := ull.Count()
			error := math.Abs(float64(estimate)-float64(tt.n)) / float64(tt.n)

			if error > tt.tolerance {
				t.Errorf("estimate = %d, actual = %d, error = %.2f%%, tolerance = %.2f%%",
					estimate, tt.n, error*100, tt.tolerance*100)
			}
		})
	}
}

func TestDuplicates(t *testing.T) {
	ull := MustNew(14)
	rng := rand.New(rand.NewPCG(42, 12345))

	// Generate 100 unique values
	values := make([]uint64, 100)
	for i := range values {
		values[i] = rng.Uint64()
	}

	// Add each value 10 times
	for range 10 {
		for _, v := range values {
			ull.Add(v)
		}
	}

	estimate := ull.Count()
	error := math.Abs(float64(estimate)-100) / 100

	if error > 0.15 {
		t.Errorf("estimate = %d, expected ~100, error = %.2f%%", estimate, error*100)
	}
}

func TestMerge(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		ull1 := MustNew(14)
		ull2 := MustNew(14)
		rng := rand.New(rand.NewPCG(42, 12345))

		// Add 500 unique values to each
		for range 500 {
			ull1.Add(rng.Uint64())
		}
		for range 500 {
			ull2.Add(rng.Uint64())
		}

		// Merge
		if err := ull1.Merge(ull2); err != nil {
			t.Fatalf("merge error: %v", err)
		}

		estimate := ull1.Count()
		// Should be close to 1000
		error := math.Abs(float64(estimate)-1000) / 1000
		if error > 0.05 {
			t.Errorf("merged estimate = %d, expected ~1000, error = %.2f%%", estimate, error*100)
		}
	})

	t.Run("different_precision", func(t *testing.T) {
		ull1 := MustNew(14)
		ull2 := MustNew(12)

		err := ull1.Merge(ull2)
		if err == nil {
			t.Error("expected error when merging different precisions")
		}
	})

	t.Run("overlapping", func(t *testing.T) {
		ull1 := MustNew(14)
		ull2 := MustNew(14)
		rng := rand.New(rand.NewPCG(42, 12345))

		// Generate shared values
		sharedValues := make([]uint64, 500)
		for i := range sharedValues {
			sharedValues[i] = rng.Uint64()
		}

		// Add shared values to both
		for _, v := range sharedValues {
			ull1.Add(v)
			ull2.Add(v)
		}

		// Add unique values to each
		for range 250 {
			ull1.Add(rng.Uint64())
		}
		for range 250 {
			ull2.Add(rng.Uint64())
		}

		if err := ull1.Merge(ull2); err != nil {
			t.Fatalf("merge error: %v", err)
		}

		estimate := ull1.Count()
		// Should be close to 1000 (500 shared + 250 + 250)
		error := math.Abs(float64(estimate)-1000) / 1000
		if error > 0.05 {
			t.Errorf("merged estimate = %d, expected ~1000, error = %.2f%%", estimate, error*100)
		}
	})
}

func TestClone(t *testing.T) {
	ull := MustNew(14)
	rng := rand.New(rand.NewPCG(42, 12345))

	for range 1000 {
		ull.Add(rng.Uint64())
	}

	clone := ull.Clone()

	// Verify clone has same count
	if ull.Count() != clone.Count() {
		t.Errorf("clone count = %d, original count = %d", clone.Count(), ull.Count())
	}

	// Verify modifications to clone don't affect original
	for range 1000 {
		clone.Add(rng.Uint64())
	}

	if clone.Count() == ull.Count() {
		t.Error("clone modifications should not affect original")
	}
}

func TestReset(t *testing.T) {
	ull := MustNew(14)
	rng := rand.New(rand.NewPCG(42, 12345))

	for range 1000 {
		ull.Add(rng.Uint64())
	}

	if ull.Count() == 0 {
		t.Error("count should be > 0 before reset")
	}

	ull.Reset()

	if ull.Count() != 0 {
		t.Errorf("count after reset = %d, want 0", ull.Count())
	}
}

func TestMarshalUnmarshal(t *testing.T) {
	ull := MustNew(14)
	rng := rand.New(rand.NewPCG(42, 12345))

	for range 1000 {
		ull.Add(rng.Uint64())
	}

	data, err := ull.MarshalBinary()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	// Expected size: 1 byte for precision + 2^14 bytes for registers
	expectedSize := 1 + (1 << 14)
	if len(data) != expectedSize {
		t.Errorf("marshal size = %d, want %d", len(data), expectedSize)
	}

	ull2 := &UltraLogLog{}
	if err := ull2.UnmarshalBinary(data); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if ull.Precision() != ull2.Precision() {
		t.Errorf("precision mismatch: %d vs %d", ull.Precision(), ull2.Precision())
	}

	if ull.Count() != ull2.Count() {
		t.Errorf("count mismatch: %d vs %d", ull.Count(), ull2.Count())
	}
}

func TestUnmarshalErrors(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"empty", []byte{}},
		{"invalid_precision_low", []byte{3}},
		{"invalid_precision_high", []byte{19}},
		{"wrong_length", []byte{14, 1, 2, 3}}, // precision 14 needs 16385 bytes
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ull := &UltraLogLog{}
			if err := ull.UnmarshalBinary(tt.data); err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestRegisterValue(t *testing.T) {
	tests := []struct {
		old      uint8
		nlz      int
		expected uint8
	}{
		{old: 0, nlz: 0, expected: 1 << 2},
		{old: 0, nlz: 1, expected: 2 << 2},
		{old: 0, nlz: 2, expected: 3 << 2},
		{old: 1 << 2, nlz: 1, expected: 2<<2 | 0b10},
		{old: 1<<2 | 0b01, nlz: 1, expected: 2<<2 | 0b10},
		{old: 1<<2 | 0b10, nlz: 1, expected: 2<<2 | 0b11},
		{old: 1<<2 | 0b10, nlz: 0, expected: 1<<2 | 0b10},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%b-%d", tt.old, tt.nlz), func(t *testing.T) {
			result := registerValue(tt.old, tt.nlz)
			if result != tt.expected {
				t.Errorf("registerValue(%.8b, %d) = %.8b, want %.8b", tt.old, tt.nlz, result, tt.expected)
			}
		})
	}
}

func TestPack(t *testing.T) {
	tests := []struct {
		reg      uint64
		expected uint8
	}{
		// {reg: 0x0000_0000_0000_0000, expected: 0b0000_0100}, // TODO: is this right??
		{reg: 0x8000_0000_0000_0000, expected: 62 << 2},
		{reg: 0xC000_0000_0000_0000, expected: 62<<2 | 0b10},
		{reg: 0x4000_0000_0000_0000, expected: 61 << 2},
		{reg: 0x2000_0000_0000_0000, expected: 60 << 2},
		{reg: 0x1000_0000_0000_0000, expected: 59 << 2},
		{reg: 0x0800_0000_0000_0000, expected: 58 << 2},
		{reg: 0x0400_0000_0000_0000, expected: 57 << 2},
		{reg: 0x0500_0000_0000_0000, expected: 57<<2 | 0b01},
		{reg: 0x0600_0000_0000_0000, expected: 57<<2 | 0b10},
		{reg: 0x0700_0000_0000_0000, expected: 57<<2 | 0b11},
		{reg: 0x0000_0000_0000_0100, expected: 7 << 2},
		{reg: 0x0000_0000_0000_0004, expected: 1 << 2},
		{reg: 0x0000_0000_0000_0008, expected: 2 << 2},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%X", tt.reg), func(t *testing.T) {
			result := pack(tt.reg)
			if result != tt.expected {
				t.Errorf("pack(%X) = %.8b, want %.8b", tt.reg, result, tt.expected)
			}
		})
	}
}

func TestUnpack(t *testing.T) {
	tests := []struct {
		packed   uint8
		expected uint64
	}{
		{packed: 0b0000_0000, expected: 0},
		{packed: 62 << 2, expected: 0x8000_0000_0000_0000},
		{packed: 62<<2 | 0b10, expected: 0xC000_0000_0000_0000},
		{packed: 61 << 2, expected: 0x4000_0000_0000_0000},
		{packed: 60 << 2, expected: 0x2000_0000_0000_0000},
		{packed: 59 << 2, expected: 0x1000_0000_0000_0000},
		{packed: 58 << 2, expected: 0x0800_0000_0000_0000},
		{packed: 57 << 2, expected: 0x0400_0000_0000_0000},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%b", tt.packed), func(t *testing.T) {
			result := unpack(tt.packed)
			if result != tt.expected {
				t.Errorf("unpack(%.8b) = %X, want %X", tt.packed, result, tt.expected)
			}
		})
	}
}
