//go:build amd64 && !js

package main

func getWinsAndLossesAVX2(b, e uint64) (w, l uint64)
func pdep(src, mask uint64) uint64
func selectBestEdgeAVX2(qs []float32, us []float32, coeff float32) int
