package weiqi

import (
	"testing"
)

func TestAll(t *testing.T) {
	t.Log(remoteSgf("http://duiyi.sina.com.cn/cgibo/20164/700e13-04292.sgf", "gb18030"))
}
