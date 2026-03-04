//go:build amd64

package ull

// mergeRegistersASM merges src registers into dst using AVX2 vectorized max operations.
// Both slices must have the same length.
//
//go:noescape
func mergeRegistersASM(dst, src *byte, length int)

// mergeRegisters merges src into dst by taking element-wise maximum.
// On AMD64 with AVX2, this uses vectorized instructions for better performance.
func mergeRegisters(dst, src []byte) {
	if len(dst) == 0 {
		return
	}
	mergeRegistersASM(&dst[0], &src[0], len(dst))
}
