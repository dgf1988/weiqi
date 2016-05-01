package weiqi

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"time"
	"net/http"
	"regexp"
	"os"
	"path/filepath"
)

const (
	c_shortDate   string = "2006年1月2日"
	c_longDate    string = "2006年01月02日"
	c_stdDate     string = "2006-01-02"
	c_stdDatetime string = "2006-01-02 15:04"
)

//ParseDate 解析日期字符串
func parseDate(dateStr string) (time.Time, error) {
	var (
		date time.Time
		err  error
	)
	date, err = time.Parse(c_stdDate, dateStr)
	if err != nil {
		date, err = time.Parse(c_longDate, dateStr)
		if err != nil {
			date, err = time.Parse(c_shortDate, dateStr)
			if err != nil {
				return time.Time{}, err
			}
		}
	}
	return date, err
}

func md5String(data string) string {
	hashMd5 := md5.New()
	io.WriteString(hashMd5, data)
	return fmt.Sprintf("%x", hashMd5.Sum(nil))
}

func md5Bytes(data []byte) string {
	hashmd5 := md5.New()
	if _, err := hashmd5.Write(data); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", hashmd5.Sum(nil))
}


func ipFromRequest(r *http.Request) string {
	var ip = r.Header.Get("x-forwarded-for")
	if ip == "" {
		ip = r.RemoteAddr
		return regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}|::1`).FindString(ip)
	}
	return ip
}

func atoi(num string) int {
	n, err := strconv.Atoi(num)
	if err != nil {
		return 0
	}
	return n
}

func atoi64(num string) int64 {
	n, err := strconv.ParseInt(num, 10, 64)
	if err != nil {
		return 0
	}
	return n
}


func mkdirIfNotExist(pathname string) error {
	var dirinfo os.FileInfo
	var err error
	if dirinfo, err = os.Stat(pathname); err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(pathname, 0666)
		} else {
			return err
		}
	} else if !dirinfo.IsDir() {
		return os.MkdirAll(pathname, 0666)
	}
	return nil
}

func addFile(filename string, data []byte) error {
	var pathname = filepath.Dir(filename)
	var err error
	var pathinfo os.FileInfo
	if pathinfo, err = os.Stat(pathname); err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(pathname, 0666); err != nil {
				return err
			}
		} else {
			return err
		}
	} else if !pathinfo.IsDir() {
		if err = os.MkdirAll(pathname, 0666); err != nil {
			return err
		}
	}
	var savef *os.File
	if savef, err = os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666); err != nil {
		return err
	} else {
		defer savef.Close()
		if _, err = savef.Write(data); err != nil {
			return err
		}
		return nil
	}
}

func removeFile(filename string) error {
	return os.Remove(filename)
}