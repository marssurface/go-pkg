package md5

import (
	"fmt"
	"testing"
)

var testData = map[string]string{
	"123456": "e10adc3949ba59abbe56e057f20f883e",
	"admin":  "21232f297a57a5a743894a0e4a801fc3",
}

func TestMd5(t *testing.T) {

	for k, v := range testData {
		if v != Md5(k) {
			t.Errorf("加密 %s 没有通过", k)
		}
	}
}


func TestMd5V2(t *testing.T) {

	for k, v := range testData {
		if v != Md5V2(k) {
			t.Errorf("加密 %s 没有通过", k)
		}
	}
}



func TestMd5V3(t *testing.T) {

	for k, v := range testData {
		if v != Md5V3(k) {
			fmt.Println(Md5V3(k))
			t.Errorf("加密 %s 没有通过", k)
		}
	}
}
