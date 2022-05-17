package qmd5

import (
	"crypto/md5"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQMD5(t *testing.T) {
	a := []byte{0x01, 0x02, 0x03}
	ma := QMD5(a)

	b := []byte{0x01, 0x02, 0x03}
	mb := QMD5(b)

	if !assert.Equal(t, ma, mb, "md5 should be same") {
		return
	}
}

func normalMD5(in []byte) []byte {
	h := md5.New()
	h.Write(in)
	return h.Sum(nil)
}

func BenchmarkQMD5(b *testing.B) {
	a := []byte{0x01, 0x02, 0x03}
	for i := 0; i < b.N; i++ {
		_ = QMD5(a)
	}
}

func BenchmarkNormalMD5(b *testing.B) {
	a := []byte{0x01, 0x02, 0x03}
	for i := 0; i < b.N; i++ {
		_ = normalMD5(a)
	}
}
