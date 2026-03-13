// Package swarlower demonstrates SWAR (SIMD Within A Register) technique for
// ASCII lowercasing — converting 'A'-'Z' to 'a'-'z' eight bytes at a time.
//
// The key insight: uppercase and lowercase ASCII letters differ only in bit 5
// (0x20).  We detect uppercase bytes using carry-safe range comparison, then
// shift the detection mask from bit 7 to bit 5 and XOR to flip case.
//
// SWAR reference: https://programming.sirrida.de/swar.html
package swarlower

import (
	"encoding/binary"
	"unsafe"
)

// LowerSlow lowercases ASCII bytes one at a time.
func LowerSlow(s string) string {
	buf := []byte(s)
	for i, b := range buf {
		if b >= 'A' && b <= 'Z' {
			buf[i] = b + 0x20
		}
	}
	return string(buf)
}

// getUint64LE reads a little-endian uint64 from the string's underlying bytes.
func getUint64LE(p *byte, off int) uint64 {
	return binary.LittleEndian.Uint64(unsafe.Slice(p, off+8)[off:])
}

// putUint64LE writes a little-endian uint64 into buf at offset.
func putUint64LE(buf []byte, off int, v uint64) {
	binary.LittleEndian.PutUint64(buf[off:], v)
}

// upperMask returns a mask with bit 7 set in each byte that is an ASCII
// uppercase letter ('A'-'Z').  All other bits are zero.
//
// The sentinel technique: we set bit 7 in every byte before subtracting so
// that borrows cannot propagate across byte lanes.
//
//	t = word | 0x8080808080808080
//	geA  = (t - 'A'broadcast) & 0x80…   // bit 7 set where byte >= 'A'
//	geBr = (t - '['broadcast) & 0x80…   // bit 7 set where byte >= '[' (i.e. > 'Z')
//	mask = geA &^ geBr                  // bit 7 set where 'A' <= byte <= 'Z'
func upperMask(word uint64) uint64 {
	const (
		broadcastA  = 0x4141414141414141 // 'A'
		broadcastBr = 0x5b5b5b5b5b5b5b5b // 'Z' + 1 = '['
		mask7       = 0x8080808080808080
	)
	// Non-ASCII bytes (>= 0x80) already have bit 7 set, which breaks the
	// sentinel technique.  Identify ASCII bytes first and exclude the rest.
	ascii := ^word & mask7 // bit 7 set only where the original byte < 0x80
	t := word | mask7
	geA := (t - broadcastA) & mask7
	geBr := (t - broadcastBr) & mask7
	return geA &^ geBr & ascii
}

// lowerWord lowercases all ASCII uppercase bytes in a uint64 word.
//
//	mask  = upperMask(word)   // 0x80 per uppercase byte
//	shift = mask >> 2         // 0x80 >> 2 = 0x20 per uppercase byte
//	result = word ^ shift     // flip bit 5 to convert case
func lowerWord(word uint64) uint64 {
	mask := upperMask(word)
	return word ^ (mask >> 2) // 0x80 → 0x20: bit 5 flip
}

// LowerSwar lowercases ASCII bytes eight at a time using SWAR.
// Non-ASCII bytes are left untouched.
func LowerSwar(s string) string {
	length := len(s)
	buf := make([]byte, length)
	i := 0

	if length >= 8 {
		p := unsafe.StringData(s)
		for ; i <= length-8; i += 8 {
			word := getUint64LE(p, i)
			putUint64LE(buf, i, lowerWord(word))
		}
	}

	// Handle remaining bytes.
	for ; i < length; i++ {
		b := s[i]
		if b >= 'A' && b <= 'Z' {
			b += 0x20
		}
		buf[i] = b
	}

	return string(buf)
}
