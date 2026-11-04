package trusted

import "testing"

func TestMemoryExhaustion(t *testing.T) {
	var chunks [][]byte
	for i := 1; ; i++ {
		b := make([]byte, 10<<20)
		for j := range b {
			b[j] = 1
		}
		chunks = append(chunks, b)
		if i%10 == 0 {
			t.Logf("allocated %d MiB", i*10)
		}
	}
}
