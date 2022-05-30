package minappCode2Session

import (
	"testing"

	_ "github.com/lib/pq"
)

func TestA(t *testing.T) {
	info := GET("wxb71c87a341a6eda7", "0c7b48c636ed0669a1a00ececdbcd22d", "0812dxFa1VuKfD0UysFa1XWLZj42dxF1")
	t.Logf("info=%#+v\n", info)
}
