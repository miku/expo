package kv

func customFunc1(s string, size int) int {
	if len(s) == 0 {
		return 0
	}
	return (int(s[0]-'a') % size)
}

func customFunc2(s string, size int) int {
	if len(s) == 0 {
		return 0
	}
	v := 0
	for i := 0; i < len(s); i++ {
		v = v + int(s[i]-'a')*i
	}
	return v % size
}
