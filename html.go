package weiqi

import (
	"fmt"
	"html/template"
	"io"
)

const (
	c_HtmlBasePath = "html/"
	c_HtmlSuffix = "html"
	c_HtmlDefClsname = "layout"
	c_HtmlDefFilename = "default"
)

func getHtmlFullname(clsname, filename string) string {
	return fmt.Sprintf("%s%s%s/%s.%s", config.BasePath, c_HtmlBasePath, clsname, filename, c_HtmlSuffix)
}

type Html struct {
	ClsName  string
	FileName string
	Childs   []*Html
}

func (h *Html) Append(htmls ...*Html) *Html {
	h.Childs = append(h.Childs, htmls...)
	return h
}

func (h *Html) Fullname() string {
	return getHtmlFullname(h.ClsName, h.FileName)
}

func (h *Html) AllFullname() []string {
	all_name := make([]string, 0)
	all_name = append(all_name, h.Fullname())
	for i := range h.Childs {
		all_name = append(all_name, h.Childs[i].AllFullname()...)
	}
	return all_name
}

func (h *Html) Execute(out io.Writer, datamap interface{}, funcmap template.FuncMap) error {
	var (
		t   = template.New("")
		err error
	)
	if funcmap != nil {
		t, err = t.Funcs(funcmap).ParseFiles(h.AllFullname()...)
	} else {
		t, err = t.ParseFiles(h.AllFullname()...)
	}
	if err != nil {
		return err
	}
	return t.ExecuteTemplate(out, h.ClsName, datamap)
}

//
func newHtml(clsname, filename string) *Html {
	if clsname == "" {
		clsname = c_HtmlDefClsname
	}
	if filename == "" {
		filename = c_HtmlDefFilename
	}
	return &Html{clsname, filename, make([]*Html, 0)}
}

//布局
func newHtmlLayout(filename string) *Html {
	return newHtml("layout", filename)
}

func defHtmlLayout() *Html {
	return newHtmlLayout("")
}

//html头
func newHtmlHead(filename string) *Html {
	return newHtml("head", filename)
}

func defHtmlHead() *Html {
	return newHtmlHead("")
}

//页面头
func newHtmlHeader(filename string) *Html {
	return newHtml("header", filename)
}

func defHtmlHeader() *Html {
	return newHtmlHeader("")
}

//页面脚
func newHtmlFooter(filename string) *Html {
	return newHtml("footer", filename)
}

func defHtmlFooter() *Html {
	return newHtmlFooter("")
}

//页面内容
func newHtmlContent(filename string) *Html {
	return newHtml("content", filename)
}

func defHtmlContent() *Html {
	return newHtmlContent("")
}
