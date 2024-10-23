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
