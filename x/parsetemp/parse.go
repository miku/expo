// https://www.reddit.com/r/golang/comments/xv9yyv/strconvparsefloat_faster_altrernatives/
package parsetemp

import (
	"math/big"
	"strconv"
)

func ParseTempFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
}

func ParseTempBigFloat(s string) float64 {
	ff, _, err := big.ParseFloat(s, 10, 64, big.ToNearestEven)
	if err != nil {
		panic(err)
	}
	f, _ := ff.Float64()
	return f
}

// ParseIntoInt parses a string containing a fixed precision (e.g. "10.2")
// float into an int (e.g. 102).
func ParseTempToInt(p []byte) int {
	var (
		result int
		pos    = 1 // exp
		digit  byte
	)
	for i := len(p) - 1; i > -1; i-- {
		if p[i] == '.' {
			continue
		} else if p[i] == '-' {
			return -result
		} else {
			digit = p[i] - '0'
			result = result + int(digit)*pos
			pos = 10 * pos
		}
	}
	return result
}

// ParseIntoIntShift parses a string containing a fixed precision (e.g. "10.2")
// float into an int (e.g. 102).
func ParseTempToIntShift(p []byte) int {
	var (
		result int
		pos    = 1 // exp
		digit  byte
	)
	for i := len(p) - 1; i > -1; i-- {
		if p[i] == '.' {
			continue
		} else if p[i] == '-' {
			return -result
		} else {
			digit = p[i] - '0'
			result = result + int(digit)*pos
			pos = (pos << 3) + (pos << 1) // pos = 10 * pos
		}
	}
	return result
}

// ParseNumber reads decimal number that matches "^-?[0-9]{1,2}[.][0-9]" pattern,
// e.g.: -12.3, -3.4, 5.6, 78.9 and return the value*10, i.e. -123, -34, 56, 789.
// From: https://github.com/gunnarmorling/1brc/blob/db064194be375edc02d6dbcd21268ad40f7e2869/src/main/go/AlexanderYastrebov/calc.go#L261C1-L283C2
func ParseNumber(data []byte) int64 {
	negative := data[0] == '-'
	if negative {
		data = data[1:]
	}
	var result int64
	switch len(data) {
	// 1.2
	case 3:
		result = int64(data[0])*10 + int64(data[2]) - '0'*(10+1)
		// 12.3
	case 4:
		result = int64(data[0])*100 + int64(data[1])*10 + int64(data[3]) - '0'*(100+10+1)
	}
	if negative {
		return -result
	}
	return result
}
