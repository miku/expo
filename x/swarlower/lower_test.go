package swarlower

import (
	"fmt"
	"strings"
	"testing"
)

var testCases = []struct {
	name string
	in   string
	want string
}{
	{"empty", "", ""},
	{"already lower", "hello world", "hello world"},
	{"all upper", "HELLO WORLD", "hello world"},
	{"mixed case", "GoLang SWAR", "golang swar"},
	{"digits untouched", "ABC123DEF", "abc123def"},
	{"punctuation", "Hello, World! (2024)", "hello, world! (2024)"},
	{"boundaries", "@AZ[`az{", "@az[`az{"},                 // chars adjacent to A-Z and a-z
	{"unicode passthrough", "Ünîcödé ABC", "Ünîcödé abc"},  // non-ASCII bytes unchanged
	{"short", "Hi", "hi"},                                  // shorter than 8 bytes
	{"exactly 8", "ABCDEFGH", "abcdefgh"},                  // exactly one word
	{"exactly 16", "ABCDEFGHIJKLMNOP", "abcdefghijklmnop"}, // two full words
	{"17 bytes", "ABCDEFGHIJKLMNOPQ", "abcdefghijklmnopq"}, // two words + tail
	{"long", strings.Repeat("Hello World! ", 100), strings.Repeat("hello world! ", 100)},
}

func TestLowerSlow(t *testing.T) {
	for _, tc := range testCases {
		got := LowerSlow(tc.in)
		if got != tc.want {
			t.Errorf("%s: LowerSlow(%q) = %q, want %q", tc.name, tc.in, got, tc.want)
		}
	}
}

func TestLowerSwar(t *testing.T) {
	for _, tc := range testCases {
		got := LowerSwar(tc.in)
		if got != tc.want {
			t.Errorf("%s: LowerSwar(%q) = %q, want %q", tc.name, tc.in, got, tc.want)
		}
	}
}

func TestLowerSwarMatchesSlow(t *testing.T) {
	// Exhaustive single-byte test: every possible byte value.
	for b := 0; b < 256; b++ {
		s := string([]byte{byte(b)})
		slow := LowerSlow(s)
		swar := LowerSwar(s)
		if slow != swar {
			t.Errorf("byte 0x%02x: slow=%q swar=%q", b, slow, swar)
		}
	}
}

func BenchmarkLowerSlow(b *testing.B) {
	s := "The Quick Brown Fox Jumps Over 13 Lazy Dogs. Year: 2024."
	b.SetBytes(int64(len(s)))
	for i := 0; i < b.N; i++ {
		LowerSlow(s)
	}
}

func BenchmarkLowerSwar(b *testing.B) {
	s := "The Quick Brown Fox Jumps Over 13 Lazy Dogs. Year: 2024."
	b.SetBytes(int64(len(s)))
	for i := 0; i < b.N; i++ {
		LowerSwar(s)
	}
}

func ExampleLowerSwar() {
	fmt.Println(LowerSwar("Hello, WORLD! 2024"))
	// Output: hello, world! 2024
}

// $ make bench
// go test -bench=. -benchmem -benchtime=5s
// goos: linux
// goarch: amd64
// pkg: github.com/miku/expo/x/swarlower
// cpu: Intel(R) Core(TM) Ultra 7 258V
// BenchmarkLowerSlow-8    32514576               177.1 ns/op       316.13 MB/s         128 B/op          2 allocs/op
// BenchmarkLowerSwar-8    45608842               126.1 ns/op       444.19 MB/s         128 B/op          2 allocs/op
// PASS
// ok      github.com/miku/expo/x/swarlower        11.847s
