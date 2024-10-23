package semi

import (
	"encoding/binary"
	"math/bits"
)

const BroadcastSemicolon = 0x3B3B3B3B3B3B3B3B
const Broadcast0x01 = 0x0101010101010101
const Broadcast0x80 = 0x8080808080808080

func encode(w []byte) uint64 {
	return binary.NativeEndian.Uint64(w)
}

func semicolonMatchBits(word uint64) uint64 {
	diff := word ^ BroadcastSemicolon
	return (diff - Broadcast0x01) & (^diff & Broadcast0x80)
}

func calcNameLen(b uint64) int {
	// "Dividing by eight gives the position of the byte"
	// https://richardstartin.github.io/posts/finding-bytes#finding-null-terminators-without-branches
	return (bits.TrailingZeros64(b) >> 3)
}

// The method maskWord() takes a long containing 8 bytes of input data and
// zeroes out all the bytes beyond the semicolon.
func maskWord(word, matchBits uint64) uint64 {
	mask := matchBits ^ (matchBits - 1)
	return word & mask
}
