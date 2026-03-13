package swardigits

import (
	"fmt"
	"testing"
)

var testCases = []struct {
	name string
	s    string
	want int
}{
	{"empty", "", 0},
	{"no digits", "hello world", 0},
	{"all digits", "1234567890", 10},
	{"mixed", "a1b2c3d4e5", 5},
	{"long text", "The year is 2024 and the temperature is 23.5 degrees", 7},
	{"unicode", "日本語123abc456", 6},
	{"tricky non-ascii", "°±²³´µ5", 1}, // 0xB0-0xB5 would false-match without ASCII guard
}

func TestCountDigits(t *testing.T) {
	for _, tc := range testCases {
		slow := CountDigitsSlow(tc.s)
		swar := CountDigitsSwar(tc.s)
		if slow != swar {
			t.Errorf("%s: slow=%d, swar=%d", tc.name, slow, swar)
		}
		if slow != tc.want {
			t.Errorf("%s: expected %d, got %d", tc.name, tc.want, slow)
		}
	}
}

func TestExtractDigits(t *testing.T) {
	for _, tc := range testCases {
		slow := ExtractDigitsSlow(tc.s)
		fast := ExtractDigitsFast(tc.s)
		if string(slow) != string(fast) {
			t.Errorf("%s: slow=%s, fast=%s", tc.name, slow, fast)
		}
	}
}

func BenchmarkCountDigitsSlow(b *testing.B) {
	testStr := "The quick brown fox jumps over 13 lazy dogs. Year: 2024."
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CountDigitsSlow(testStr)
	}
}

func BenchmarkCountDigitsSwar(b *testing.B) {
	testStr := "The quick brown fox jumps over 13 lazy dogs. Year: 2024."
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CountDigitsSwar(testStr)
	}
}

func BenchmarkExtractDigitsFast(b *testing.B) {
	testStr := "The quick brown fox jumps over 13 lazy dogs. Year: 2024."
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ExtractDigitsFast(testStr)
	}
}

// go test -v -bench=. -benchmem ./...
// goos: linux
// goarch: amd64
// pkg: github.com/miku/expo/x/swardigits
// cpu: Intel(R) Core(TM) i7-8565U CPU @ 1.80GHz
// BenchmarkCountDigitsSlow-8             50000000                23.8 ns/op
// BenchmarkCountDigitsSwar-8            1000000000                 0.63 ns/op
// BenchmarkExtractDigitsFast-8          2000000000                 0.58 ns/op
// PASS

func ExampleCountDigitsSwar() {
	s := "The year is 2024 and the temperature is 23.5 degrees"
	count := CountDigitsSwar(s)
	fmt.Printf("Found %d digits\n", count)
	// Output: Found 7 digits
}

func ExampleExtractDigitsFast() {
	s := "abc123def456"
	digits := ExtractDigitsFast(s)
	fmt.Printf("Digits: %s\n", digits)
	// Output: Digits: 123456
}
