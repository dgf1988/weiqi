package weiqi

import (
	"fmt"
	"html/template"
	"strings"
	"time"
)

type Post struct {
	Id      int64
	Title   string
	Text    string
	Pstatus int64
	Pposted time.Time
	Pupdate time.Time
}

func (p *Post) HtmlText() template.HTML {
	return template.HTML(p.Text)
}

const Post_CutText_Length = 140

func cutPostText(text string, length int) string {
	if length <= 0 {
		length = Post_CutText_Length
	}
	s := []rune(text)
	if len(s) < length {
		return text
	}
	return string(s[:length])
}

func cutPostTextMany(ps []Post) {
	if ps == nil {
		return
	}
	for i := range ps {
		ps[i].Text = cutPostText(ps[i].Text, 0)
	}
}

func parseTextToHtml(text string) string {
	lines := strings.Split(text, "\n")
	ret := make([]string, 0)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "<p") && strings.HasSuffix(line, "</p>") {
			ret = append(ret, line)
		} else {
			ret = append(ret, fmt.Sprint("<p>", line, "</p>"))
		}
	}
	return strings.Join(ret, "\n")
}
