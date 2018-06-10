package cuckoofilters

import (
	"crypto/sha1"
	"encoding/binary"
)

func hash(data []byte) uint64 {
	s := sha1.Sum(data)
	return binary.LittleEndian.Uint64(s[:])
}

func nearestPowerOfTwo(val uint32) uint32 {
	for i := uint32(0); i < 32; i++ {
		if pow := uint32(1) << i; pow >= val {
			return pow
		}
	}
	panic(errNearestPowerOfTwoStr)
}
