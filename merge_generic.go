package ull

// mergeRegisters
func mergeRegisters(dst, src []byte) {
	for i, val := range src {
		if val == 0 {
			continue
		}
		dstVal := dst[i]
		dst[i] = pack(unpack(val) | unpack(dstVal))
	}
}
