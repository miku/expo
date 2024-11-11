package semi

import (
	"testing"
)

func TestSemicolonMatchBits(t *testing.T) {
	var cases = []struct {
		B []byte
	}{
		{[]byte("1;222222")},
		{[]byte("hello;wo")},
		{[]byte("hellowor")},
	}
	for _, c := range cases {
		t.Log(string(c.B))
		t.Log(c.B)
		w := encode(c.B)
		t.Log(w)
		t.Logf("%064b (word)", w)
		matchBits := semicolonMatchBits(w)
		t.Log(matchBits)
		t.Logf("%064b (matchBits)", matchBits)
		if matchBits == 0 {
			t.Logf("no ';' match")
			continue
		}
		L := calcNameLen(matchBits)
		t.Logf("%064b (nameLen)", L)
		masked := maskWord(w, matchBits)
		t.Log(masked)
		t.Logf("%064b (masked)", masked)
		t.Logf("----")
	}
}

// go test -v .
// === RUN   TestSemicolonMatchBits
//     semi_test.go:16: 1;222222
//     semi_test.go:17: [49 59 50 50 50 50 50 50]
//     semi_test.go:19: 3617008641903835953
//     semi_test.go:20: 0011001000110010001100100011001000110010001100100011101100110001 (word)
//     semi_test.go:22: 32768
//     semi_test.go:23: 0000000000000000000000000000000000000000000000001000000000000000 (matchBits)
//     semi_test.go:29: 0000000000000000000000000000000000000000000000000000000000000001 (nameLen)
//     semi_test.go:31: 15153
//     semi_test.go:32: 0000000000000000000000000000000000000000000000000011101100110001 (masked)
//     semi_test.go:33: ----
//     semi_test.go:16: hello;wo
//     semi_test.go:17: [104 101 108 108 111 59 119 111]
//     semi_test.go:19: 8031953810185020776
//     semi_test.go:20: 0110111101110111001110110110111101101100011011000110010101101000 (word)
//     semi_test.go:22: 140737488355328
//     semi_test.go:23: 0000000000000000100000000000000000000000000000000000000000000000 (matchBits)
//     semi_test.go:29: 0000000000000000000000000000000000000000000000000000000000000101 (nameLen)
//     semi_test.go:31: 65349746451816
//     semi_test.go:32: 0000000000000000001110110110111101101100011011000110010101101000 (masked)
//     semi_test.go:33: ----
//     semi_test.go:16: hellowor
//     semi_test.go:17: [104 101 108 108 111 119 111 114]
//     semi_test.go:19: 8245940763182785896
//     semi_test.go:20: 0111001001101111011101110110111101101100011011000110010101101000 (word)
//     semi_test.go:22: 0
//     semi_test.go:23: 0000000000000000000000000000000000000000000000000000000000000000 (matchBits)
//     semi_test.go:25: no ';' match
// --- PASS: TestSemicolonMatchBits (0.00s)
// PASS
