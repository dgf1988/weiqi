package weiqi

import (
	"strings"
)

//导航结构
type NavItem struct {
	Title string
	Url   string
}

//默认导航
func defNavItems() []NavItem {
	return []NavItem{NavItem{"首页", "/"}, NavItem{"文章", "/post/"}, NavItem{"棋谱", "/sgf/"}, NavItem{"棋手", "/player/"}}
}

//用户管理导航
func userNavItems() []NavItem {
	return []NavItem{NavItem{"首页", "/"}, NavItem{"用户", "/user"}, NavItem{"文章", "/user/post/"}, NavItem{"棋谱", "/user/sgf/"}, NavItem{"棋手", "/user/player/"}, NavItem{"图片", "/user/img/"}}
}

const (
	c_lengthFy = 11
)

//翻页
type IndexPageItem struct {
	IsCurrent bool
	Number    int
}

type IndexPages struct {
	First *IndexPageItem
	Last *IndexPageItem
	Indexs []IndexPageItem
}

func newIndexPages(currnet, last int) *IndexPages {
	if currnet > last || last < 1{
		return nil
	}
	if last == 1 {
		return &IndexPages{&IndexPageItem{true, 1}, nil, nil}
	}
	if last == 2 {
		return &IndexPages{&IndexPageItem{currnet == 1, 1}, &IndexPageItem{currnet == 2, 2}, nil}
	}
	if last <= c_lengthFy {
		return &IndexPages{&IndexPageItem{currnet == 1, 1}, &IndexPageItem{currnet == last, last}, makeIndexPageItems(currnet, 2, -1 + last)}
	}
	var harf int = c_lengthFy/2
	if currnet <= harf {
		return &IndexPages{&IndexPageItem{currnet == 1, 1}, &IndexPageItem{currnet == last, last}, makeIndexPageItems(currnet, 2, -1 + c_lengthFy)}
	}
	if currnet >= last - harf {
		return &IndexPages{&IndexPageItem{currnet == 1, 1}, &IndexPageItem{currnet == last, last}, makeIndexPageItems(currnet, 2 + last - c_lengthFy, -1 + last)}
	}
	return &IndexPages{&IndexPageItem{currnet == 1, 1}, &IndexPageItem{currnet == last, last}, makeIndexPageItems(currnet, currnet - harf + 1, currnet + harf - 1)}
}

func makeIndexPageItems(current, beg, end int) []IndexPageItem {
	indexpages := make([]IndexPageItem, 0)
	for i := beg ; i <= end ; i++ {
		indexpages = append(indexpages, IndexPageItem{i==current,i})
	}
	return indexpages
}

//Head 页面布局使用的Html头数据结构
type Head struct {
	Title    string
	Desc     string
	Keywords []string
}

//StrKeywords 用来在模板上直接输出字符串
func (h Head) StrKeywords() string {
	return strings.Join(h.Keywords, ",")
}

func defHead() *Head {
	return &Head{
		config.SiteTitle,
		config.SiteDesc,
		config.SiteKeywords,
	}
}

//页面头结构
type Header struct {
	Title string
	Navs  []NavItem
}

func defHeader() *Header {
	return &Header{
		config.SiteTitle,
		defNavItems(),
	}
}

func userHeader() *Header {
	return &Header{
		config.SiteTitle,
		userNavItems(),
	}
}

//页面脚结构
type Footer struct {
	AuthorName  string
	AuthorURL   string
	AuthorEmail string
	ICP         string
}

func defFooter() *Footer {
	return &Footer{
		config.SiteAuthorName,
		config.SiteAuthorUrl,
		config.SiteAuthorEmail,
		config.SiteICP,
	}
}

//编辑器结构
type Editor struct {
	Action string
	Msg    string
}

func newEditor(action, msg string) *Editor {
	return &Editor{
		action, msg,
	}
}

//内容数据结构
type Content map[string]interface{}

type Data struct {
	User    *User
	Head    *Head
	Header  *Header
	Content Content
	Footer  *Footer
}

func defData() *Data {
	data := &Data{}
	data.Head = defHead()
	data.Header = defHeader()
	data.Footer = defFooter()
	data.Content = make(Content)
	return data
}

func userData() *Data {
	data := &Data{}
	data.Head = defHead()
	data.Header = userHeader()
	data.Footer = defFooter()
	data.Content = make(Content)
	return data
}
