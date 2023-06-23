package utils

import (
	"crypto/md5"
	"fmt"
	"io"
)

// Md5Sign md5加密
func Md5Sign(str string) string {
	w := md5.New()
	io.WriteString(w, str)
	md5Str := fmt.Sprintf("%x", w.Sum(nil))
	return md5Str
}
