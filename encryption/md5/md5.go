package md5

/**
md5 实现
*/

import (
	"crypto"
	"crypto/md5"
	"encoding/hex"
	"io"
)

func Md5(str string) string {
	m := crypto.MD5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

func Md5V2(str string) string {
	m := md5.New()
	io.WriteString(m, str)
	return hex.EncodeToString(m.Sum(nil))
}

func Md5V3(str string) string {
	m := md5.Sum([]byte(str))
	return hex.EncodeToString(m[:])
}
