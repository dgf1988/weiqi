package main

import (
	"crypto/md5"
	"io"
	"fmt"
)

func getMd5(data string) string {
	hashMd5 := md5.New()
	io.WriteString(hashMd5, data)
	return fmt.Sprintf("%x", hashMd5.Sum(nil))
}
