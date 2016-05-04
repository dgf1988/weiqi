package weiqi

import (
	"fmt"
	"html/template"
	"strings"
	"time"
)

const (
	c_postPageSize = 10
)

type Post struct {
	Id     int64
	Title  string
	Text   string
	Status int64
	Posted time.Time
	Update time.Time
}

func (p *Post) HtmlText() template.HTML {
	return template.HTML(p.Text)
}

func (p Post) StrStatus() string {
	return formatStatus(p.Status)
}

const c_CutTextLength = 140

func cutPostText(text string, length int) string {
	if length <= 0 {
		length = c_CutTextLength
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

func listPostByStatusOrderDesc(status, take, skip int) ([]Post, error) {
	var posts = make([]Post, 0)
	if rows, err := Db.Post.Query("where post.status = ? order by post.id desc limit ?, ?", status, skip, take); err != nil {
		return nil, err
	} else {
		defer rows.Close()
		for rows.Next() {
			var post Post
			if err = rows.Struct(&post); err != nil {
				return nil, err
			} else {
				posts = append(posts, post)
			}
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}
	}
	return posts, nil
}
