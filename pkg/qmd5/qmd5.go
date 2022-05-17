package qmd5

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"hash"
	"strings"
	"sync"
)

var md5Pool = sync.Pool{
	New: func() interface{} {
		return md5.New()
	},
}

func QMD5(in []byte) []byte {
	h := md5Pool.Get().(hash.Hash)
	h.Reset()
	h.Write(in)
	defer md5Pool.Put(h)
	return h.Sum(nil)
}

func QMD5String(in []byte) string {
	return hex.EncodeToString(QMD5(in))
}

func QMD5Match(in []byte, md5 []byte) bool {
	return bytes.Equal(QMD5(in), md5)
}

func QMD5MatchString(in []byte, md5 string) bool {
	return QMD5String(in) == strings.ToLower(md5)
}
