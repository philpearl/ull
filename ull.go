// Package ull implements UltraLogLog, a cardinality estimation algorithm
// that improves upon HyperLogLog with better accuracy at the same memory cost.
//
// Reference: Otmar Ertl, "UltraLogLog: A Practical and More Space-Efficient
// Alternative to HyperLogLog for Approximate Distinct Counting" (2023)
package ull

import (
	"errors"
	"math/bits"

	"github.com/cespare/xxhash/v2"
)

// UltraLogLog is a probabilistic cardinality estimator.
// It uses a configurable number of registers (2^precision) to estimate
// the number of distinct elements added to it.
type UltraLogLog struct {
	// Each 8 bit register stores a 6-bit number and a 2-bit number.
	//
	// - the upper 6 bits store the maximum number of leading zeros (+1) seen
	// for a hash value added to the bucket represented by the register.
	//
	// - the lower 2 bits store whether we've also seen numbers of leading zeros
	// that were 1,2 or 3 less than the maximum.
	registers []uint8
	// precision is the number of bits used to determine the register index.
	precision uint8
}

// New creates a new UltraLogLog with the specified precision.
// Precision must be between 4 and 18 (inclusive).
// Higher precision uses more memory but provides better accuracy.
// Memory usage is 2^precision bytes.
//
// Recommended precision values:
//   - 10: 1KB memory, ~2.6% standard error
//   - 12: 4KB memory, ~1.3% standard error
//   - 14: 16KB memory, ~0.65% standard error
//   - 16: 64KB memory, ~0.325% standard error
func New(precision uint8) (*UltraLogLog, error) {
	if precision < 4 || precision > 18 {
		return nil, errors.New("precision must be between 4 and 18")
	}
	m := uint32(1) << precision
	return &UltraLogLog{
		registers: make([]uint8, m),
		precision: precision,
	}, nil
}

// MustNew creates a new UltraLogLog with the specified precision.
// It panics if precision is invalid.
func MustNew(precision uint8) *UltraLogLog {
	ull, err := New(precision)
	if err != nil {
		panic(err)
	}
	return ull
}

// Add adds a pre-hashed 64-bit value to the UltraLogLog.
// The hash should be a high-quality hash of the original value.
func (u *UltraLogLog) Add(hash uint64) {
	// Use the first `precision` bits to determine the bucket index
	p := u.precision
	q := 64 - p
	idx := hash >> q

	// nlz = number of leading zeros in ~(~hash << p)
	// This finds the position of the first 0 bit after the index portion
	nlz := bits.LeadingZeros64(^(^hash << p))

	old := u.registers[idx]

	// Compute bit position: (nlz + p - 1) mod 64
	bitpos := (uint(nlz) + uint(p) - 1) & 63
	hp := unpack(old)
	hp |= 1 << bitpos

	u.registers[idx] = pack(hp)
}

// pack encodes a uint64 register value into a single byte. The uint64 register
// value represents a hash value in the register bucket. It has a number of
// leading zeros, then a 1 followed by some additional bits that represent
// whether we saw entries with fewer leading zeros!
func pack(reg uint64) uint8 {
	if reg == 0 {
		return 0
	}
	nlz := bits.LeadingZeros64(reg) + 1 // 1..=64
	s := uint32(nlz & 63)
	y := uint8((reg << s) >> 62)
	return uint8(((-int32(nlz)) << 2) | int32(y))
}

// unpack decodes a register value back into a uint64 with an appropriate number
// of leading zeros, then a 1 followed by the additional bits
func unpack(val uint8) uint64 {
	if val == 0 {
		return 0
	}

	sh := ((uint32(val>>2) - 2) & 63)
	return (4 | uint64(val&3)) << sh
}

// AddBytes adds a byte slice to the UltraLogLog using xxhash.
func (u *UltraLogLog) AddBytes(data []byte) {
	u.Add(xxhash.Sum64(data))
}

// AddString adds a string to the UltraLogLog using xxhash.
func (u *UltraLogLog) AddString(s string) {
	u.Add(xxhash.Sum64String(s))
}

// Count returns the estimated cardinality of the set.
func (u *UltraLogLog) Count() uint64 {
	return u.FGRAEstimate()
}

// Merge combines another UltraLogLog into this one.
// Both must have the same precision.
func (u *UltraLogLog) Merge(other *UltraLogLog) error {
	if u.precision != other.precision {
		return errors.New("cannot merge UltraLogLogs with different precisions")
	}

	mergeRegisters(u.registers, other.registers)

	return nil
}

// Clone creates a deep copy of the UltraLogLog.
func (u *UltraLogLog) Clone() *UltraLogLog {
	clone := &UltraLogLog{
		registers: make([]uint8, len(u.registers)),
		precision: u.precision,
	}
	copy(clone.registers, u.registers)
	return clone
}

// Reset clears all registers, allowing the UltraLogLog to be reused.
func (u *UltraLogLog) Reset() {
	for i := range u.registers {
		u.registers[i] = 0
	}
}

// Precision returns the precision of the UltraLogLog.
func (u *UltraLogLog) Precision() uint8 {
	return u.precision
}

// Size returns the memory size of the registers in bytes.
func (u *UltraLogLog) Size() int {
	return len(u.registers)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (u *UltraLogLog) MarshalBinary() ([]byte, error) {
	data := make([]byte, 1+len(u.registers))
	data[0] = u.precision
	copy(data[1:], u.registers)
	return data, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (u *UltraLogLog) UnmarshalBinary(data []byte) error {
	if len(data) < 1 {
		return errors.New("invalid data: too short")
	}

	precision := data[0]
	if precision < 4 || precision > 18 {
		return errors.New("invalid precision in data")
	}

	expectedLen := 1 + (1 << precision)
	if len(data) != expectedLen {
		return errors.New("invalid data length")
	}

	u.precision = precision
	u.registers = make([]uint8, 1<<precision)
	copy(u.registers, data[1:])

	return nil
}
