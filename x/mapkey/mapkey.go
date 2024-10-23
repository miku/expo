package mapkey

import (
	"hash/fnv"
	"math/big"
	"math/bits"
	"unsafe"
)

func ByteSliceToString(p []byte) string {
	key := string(p)
	return key
}

func ByteSliceToFNV(p []byte) []byte {
	h := fnv.New32()
	return h.Sum(p)
}

func ByteSliceToIntKey(p []byte) int {
	result := 0
	for i, c := range p {
		result = i*result + int(c-'a')
	}
	return result
}

func ByteSliceToBigInt(p []byte) int {
	return int(big.NewInt(0).SetBytes(p).Uint64())
}

func ByteSliceToStringZeroAlloc(p []byte) string {
	return *(*string)(unsafe.Pointer(&p))
}

// StringToIndex taken from a custom map implementation.
func StringToIndex16(s string) (index int) {
	for i, c := range s {
		index = index + i*(37+int(c))
	}
	return index % 16384
}

// https://theartincode.stanis.me/008-djb2/
// Written by Daniel J. Bernstein (also known as djb), this simple hash function dates back to 1991.
func Djb2(s string) int64 {
	var hash int64 = 5381
	for _, c := range s {
		hash = ((hash << 5) + hash) + int64(c)
		// the above line is an optimized version of the following line:
		// hash = hash * 33 + int64(c)
		// which is easier to read and understand...
	}
	return hash
}

// FxHasher is based on a hash function used within Firefox. (Indeed, the Fx is
// short for “Firefox”.)
//
// (Are you wondering where the constant
// 0x517cc1b727220a95 comes from? 0xffff_ffff_ffff_ffff / 0x517c_c1b7_2722_0a95
// = π.)
//
// In terms of hashing quality, it is mediocre. If you run it through a hash
// quality tester it will fail a number of the tests. For example, if you hash
// any sequence of N zeroes, you get zero. And yet, for use in hash tables
// within the Rust compiler, it’s hard to beat.
func FxHasher(word uint64) uint64 {
	// https://nnethercote.github.io/2021/12/08/a-brutally-effective-hash-function-in-rust.html
	return bits.RotateLeft64(word*0x51_7c_c1_b7_27_22_0a_95, 17)
}
