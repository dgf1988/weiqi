package weiqi

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"time"
	"net/http"
	"regexp"
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

func getMd5(data string) string {
	hashMd5 := md5.New()
	io.WriteString(hashMd5, data)
	return fmt.Sprintf("%x", hashMd5.Sum(nil))
}

func getIp(r *http.Request) string {
	var ip = r.Header.Get("x-forwarded-for")
	if ip == "" {
		ip = r.RemoteAddr
		return regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`).FindString(ip)
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
