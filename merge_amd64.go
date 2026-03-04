//go:build amd64

package ull

import "golang.org/x/sys/cpu"

// mergeRegistersAVX2 merges src registers into dst using AVX2 vectorized max operations.
//
//go:noescape
func mergeRegistersAVX2(dst, src *byte, length int)

// mergeRegistersAVX512 merges src registers into dst using AVX-512 vectorized max operations.
//
//go:noescape
func mergeRegistersAVX512(dst, src *byte, length int)

// hasAVX512 is set at init time based on CPU capabilities.
var hasAVX512 = cpu.X86.HasAVX512F && cpu.X86.HasAVX512BW

// mergeRegisters merges src into dst by taking element-wise maximum.
// On AMD64, this uses AVX-512 if available, otherwise AVX2.
func mergeRegisters(dst, src []byte) {
	if len(dst) == 0 {
		return
	}
	if hasAVX512 {
		mergeRegistersAVX512(&dst[0], &src[0], len(dst))
	} else {
		mergeRegistersAVX2(&dst[0], &src[0], len(dst))
	}
}
