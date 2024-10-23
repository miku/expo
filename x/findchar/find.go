package findchar

import (
	"bytes"
	"encoding/binary"
	"math/bits"
	"strings"
)

func WithCut(s, c string) {
	_, _, _ = strings.Cut(s, c)
}

func WithSplit(s, c string) {
	_ = strings.Split(s, c)
}

func WithIndex(s, c string) {
	_ = strings.Index(s, c)
}

func WithIter(s, c string) {
	for _, v := range s {
		if v == rune(c[0]) {
			break
		}
	}
}

// SWAR related code

const BroadcastSemicolon = 0x3B3B3B3B3B3B3B3B
const Broadcast0x01 = 0x0101010101010101
const Broadcast0x80 = 0x8080808080808080

// WithSwar, here: separator is assumed a ";"
func WithSwar(b []byte) {
	for len(b) > 8 {
		word := binary.NativeEndian.Uint64(b)
		matchBits := semicolonMatchBits(word)
		if matchBits != 0 {
			_ = calcNameLen(matchBits)
			return
		}
		b = b[8:]
	}
	_ = bytes.Index(b, []byte(";"))
}

func semicolonMatchBits(word uint64) uint64 {
	diff := word ^ BroadcastSemicolon
	return (diff - Broadcast0x01) & (^diff & Broadcast0x80)
}

func calcNameLen(b uint64) int {
	return (bits.TrailingZeros64(b) >> 3)
}

func maskWord(word, matchBits uint64) uint64 {
	mask := matchBits ^ (matchBits - 1)
	return word & mask
}
