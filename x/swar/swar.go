package swar

import (
	"unsafe"
)

const asciiMask = 0x8080808080808080 // 8 bytes

func IsAsciiSlow(s string) bool {
	for i := 0; i < len(s); i++ {
		//  If we perform a bitwise AND with 0x80 the result will be non-zero
		//  if the high bit is set. 0x80 = 128
		if 0x80&s[i] > 0 {
			return false
		}
	}
	return true
}

func GetBytesUint64(bytes *byte, offset int) uint64 {
	data := *(*uint64)(unsafe.Add(unsafe.Pointer(bytes), offset))
	return data
}

func IsAsciiSwar(s string) bool {
	length := len(s)
	i := 0
	if length >= 8 {
		bytes := unsafe.StringData(s)
		for ; i < length-7; i += 8 {
			if (asciiMask & GetBytesUint64(bytes, i)) > 0 {
				return false
			}
		}
	}
	// Fall back to slower approach for the last bytes.
	for ; i < length; i++ {
		if 0x80&s[i] > 0 {
			return false
		}
	}
	return true
}

func AddUint16(a [4]uint16, b [4]uint16) [4]uint16 {
	for i := 0; i < 4; i++ {
		a[i] = a[i] + b[i]
	}
	return a
}

type Lane byte

const (
	L0 Lane = 48
	L1      = 32
	L2      = 16
	L3      = 0
)

// SetPseudoLane puts a uint16 into a uint64, using masks to specify the
// pseudo-lane.
func SetPseudoLane(n uint64, v uint16, lane Lane) uint64 {
	switch lane {
	case L0:
		return (n & 0x0000FFFFFFFFFFFF) | (uint64(v) << lane)
	case L1:
		return (n & 0xFFFF0000FFFFFFFF) | (uint64(v) << lane)
	case L2:
		return (n & 0xFFFFFFFF0000FFFF) | (uint64(v) << lane)
	case L3:
		return (n & 0xFFFFFFFFFFFF0000) | (uint64(v) << lane)
	default:
		panic("invalid lane")
	}
}

// AddUint16Swar performs addition on 4 uint16 values at once, with the
// limitation that carries are not supported. Cf.
// https://programming.sirrida.de/swar.html
func AddUint16Swar(a [4]uint16, b [4]uint16) [4]uint16 {
	var c, d, e uint64
	var result [4]uint16
	c = SetPseudoLane(c, a[0], L0)
	c = SetPseudoLane(c, a[1], L1)
	c = SetPseudoLane(c, a[2], L2)
	c = SetPseudoLane(c, a[3], L3)
	d = SetPseudoLane(d, a[0], L0)
	d = SetPseudoLane(d, a[1], L1)
	d = SetPseudoLane(d, a[2], L2)
	d = SetPseudoLane(d, a[3], L3)
	e = c + d
	result[0] = uint16(e >> 48)
	result[1] = uint16(e >> 32)
	result[2] = uint16(e >> 16)
	result[3] = uint16(e >> 0)
	return result
}
