package main

import (
	"time"
	"database/sql"
	"html/template"
	"strings"
	"fmt"
)

type P struct {
	Id int64
	Title string
	Text string
	Pstatus int64
	Pposted time.Time
	Pupdate time.Time
}

func (p *P) HtmlText() template.HTML {
	return template.HTML(p.Text)
}

func dbListPostByPage(pagesize int, page int) ([]P, error) {
	rows, err := db.Query("select * from post order by pposted desc limit ?,?", pagesize*page, pagesize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ps := make([]P, 0)
	for rows.Next() {
		var p P
		err := rows.Scan(&p.Id, &p.Title, &p.Text, &p.Pstatus, &p.Pposted, &p.Pupdate)
		if err != nil {
			return ps, err
		} else {
			ps = append(ps, p)
		}
	}
	if err = rows.Err(); err != nil {
		return ps, err
	}
	return ps, nil
}

func dbAddPost(p *P) (int64, error) {
	res, err := db.Exec("insert into post (ptitle, ptext) values (?,?)", p.Title, p.Text)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func dbGetPost(id int64) (*P, error) {
	row := db.QueryRow("select * from post where id = ?  limit 1", id)
	var p P
	err := row.Scan(&p.Id, &p.Title, &p.Text, &p.Pstatus, &p.Pposted, &p.Pupdate)
	if err != nil {
		return nil,err
	}
	return &p, nil
}

func dbUpdatePost(p *P) error {
	_, err := db.Exec("update post set ptitle = ?, ptext = ? where id = ? limit 1", p.Title, p.Text, p.Id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

func dbDeletePost(id int64) error {
	_, err := db.Exec("delete from post where id = ? limit 1", id)
	if err != nil && err != sql.ErrNoRows{
		return err
	}
	return nil
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

func cutPostTextMany(ps []P) {
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

