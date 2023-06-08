package sgr

import "testing"

func BenchmarkBtoa(b *testing.B) {
	buf := make([]byte, 0, 64)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i != b.N; i++ {
		buf = buf[0:0]
		buf = append(buf, btoa(255)...)
	}
}
