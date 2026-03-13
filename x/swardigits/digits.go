// Package swardigits demonstrates SWAR (SIMD Within A Register) technique
// for counting digit characters in a string.
//
// SWAR reference: https://programming.sirrida.de/swar.html
package swardigits

import (
	"encoding/binary"
	"math/bits"
	"unsafe"
)

// CountDigitsSlow counts digit characters one byte at a time.
func CountDigitsSlow(s string) int {
	count := 0
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			count++
		}
	}
	return count
}

// GetUint64LE gets a little-endian uint64 from bytes at offset.
func GetUint64LE(bytes *byte, offset int) uint64 {
	return binary.LittleEndian.Uint64(unsafe.Slice(bytes, offset+8)[offset:])
}

// digitMask returns a mask with bit 7 set in each byte that is an ASCII digit
// ('0'-'9'). All other bits are zero.
//
// The key SWAR trick: set bit 7 of every byte before subtracting.  This
// "sentinel" bit absorbs any borrow so that carries never propagate across
// byte lanes.  After the subtraction the sentinel tells us whether the
// original byte was >= the threshold:
//
//	(b|0x80) - threshold  →  bit 7 stays set iff b >= threshold  (for ASCII b)
//
// We test two thresholds: '0' (0x30) and ':' (0x3A, i.e. '9'+1).
// A byte is a digit when b >= '0' AND b < ':'.
//
// Note: the sentinel technique only works for bytes < 0x80 (ASCII).
// Non-ASCII bytes already have bit 7 set, which defeats the sentinel;
// for example byte 0xB3 would falsely match as a digit without the
// explicit ASCII guard below.
func digitMask(word uint64) uint64 {
	const (
		broadcastLo = 0x3030303030303030 // '0'
		broadcastHi = 0x3a3a3a3a3a3a3a3a // '9' + 1 = ':'
		mask7       = 0x8080808080808080 // bit 7 in every byte
	)

	// Non-ASCII bytes (>= 0x80) already have bit 7 set, which breaks the
	// sentinel technique.  Identify ASCII bytes first and exclude the rest.
	ascii := ^word & mask7 // bit 7 set only where the original byte < 0x80

	t := word | mask7 // set sentinel bits to block inter-byte borrows

	geLo := (t - broadcastLo) & mask7 // bit 7 set where byte >= '0'
	geHi := (t - broadcastHi) & mask7 // bit 7 set where byte >= ':'

	return geLo &^ geHi & ascii // bit 7 set where '0' <= byte <= '9' (ASCII only)
}

// CountDigitsSwar counts digit characters 8 bytes at a time using SWAR.
func CountDigitsSwar(s string) int {
	length := len(s)
	count := 0
	i := 0

	if length >= 8 {
		bytes := unsafe.StringData(s)
		for ; i < length-7; i += 8 {
			word := GetUint64LE(bytes, i)
			mask := digitMask(word)

			// Count set bits in mask (these are the digit positions)
			count += bits.OnesCount64(mask)
		}
	}

	// Handle remaining bytes
	for ; i < length; i++ {
		if s[i] >= '0' && s[i] <= '9' {
			count++
		}
	}

	return count
}

// ExtractDigitsSlow extracts digit characters one byte at a time.
func ExtractDigitsSlow(s string) []byte {
	result := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			result = append(result, s[i])
		}
	}
	return result
}

// ExtractDigitsSwar extracts all digit characters from a string using SWAR.
// Returns a slice containing only the digit bytes.
func ExtractDigitsSwar(s string) []byte {
	return ExtractDigitsFast(s)
}

// ExtractDigitsFast is an optimized version that avoids append overhead.
func ExtractDigitsFast(s string) []byte {
	length := len(s)
	result := make([]byte, length) // Max possible digits
	written := 0

	i := 0

	if length >= 8 {
		bytes := unsafe.StringData(s)
		for ; i < length-7; i += 8 {
			word := GetUint64LE(bytes, i)
			mask := digitMask(word)

			for j := 0; j < 8; j++ {
				if (mask>>(j*8))&0x80 != 0 {
					result[written] = s[i+j]
					written++
				}
			}
		}
	}

	for ; i < length; i++ {
		if s[i] >= '0' && s[i] <= '9' {
			result[written] = s[i]
			written++
		}
	}

	return result[:written]
}
