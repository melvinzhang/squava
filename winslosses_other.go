//go:build !amd64 || js

package main

func getWinsAndLossesAVX2(b, e uint64) (w, l uint64) {
	return getWinsAndLossesGo(b, e)
}

func pdep(src, mask uint64) uint64 {
	// Simple bit-by-bit pdep implementation for non-x86
	var res uint64
	for i, j := 0, 0; j < 64; j++ {
		if (mask>>j)&1 != 0 {
			if (src>>i)&1 != 0 {
				res |= 1 << j
			}
			i++
		}
	}
	return res
}
