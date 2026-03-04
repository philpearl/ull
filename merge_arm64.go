//go:build arm64

package ull

// mergeRegistersASM merges src registers into dst using vectorized max operations.
// Both slices must have the same length, which must be a multiple of 16.
//
//go:noescape
func mergeRegistersASM(dst, src *byte, length int)

// mergeRegisters merges src into dst by taking element-wise maximum.
// On ARM64, this uses NEON vectorized instructions for better performance.
func mergeRegisters(dst, src []byte) {
	if len(dst) == 0 {
		return
	}
	mergeRegistersASM(&dst[0], &src[0], len(dst))
}
