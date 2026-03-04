//go:build !arm64 && !amd64

package ull

// mergeRegisters merges src into dst by taking element-wise maximum.
// This is the pure Go implementation for non-ARM64 platforms.
func mergeRegisters(dst, src []byte) {
	for i, val := range src {
		dst[i] = max(dst[i], val)
	}
}
