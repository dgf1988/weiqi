package weiqi

import (
	"net/http"
	"github.com/dgf1988/weiqi/h"
	"mime/multipart"
	"fmt"
	"image/png"
	"image"
	"image/jpeg"
	"log"
	"image/gif"
)

func img_editor_handler(w http.ResponseWriter, r *http.Request, args []string) {
	var user = getSessionUser(r)
	if user == nil {
		h.SeeOther(w, r, "/login")
		return
	}
	var err error
	var data = defData()
	data.User = user
	data.Head.Title = "图片管理"
	data.Header.Navs = userNavItems()

	if err = defHtmlLayout().Append(defHtmlHead(), defHtmlHeader(), defHtmlFooter(), newHtmlContent("userimg")).Execute(w, data, nil); err != nil {
		logError(err.Error())
		return
	}
}

func img_upload_handler(w http.ResponseWriter, r *http.Request, args []string) {
	if getSessionUser(r) == nil {
		h.SeeOther(w, r, "/login")
		return
	}
	var err error

	if err = r.ParseMultipartForm(2<<20); err != nil {
		h.Text(w, err.Error(), 200)
		return
	}

	var f multipart.File
	var header *multipart.FileHeader
	if f, header, err = r.FormFile("file"); err != nil {
		h.Text(w, err.Error(), 200)
		return
	} else {
		var imgtype = parseImgType(header.Header.Get("content-type"))
		var imgconf image.Config
		switch imgtype {
		case c_img_png:
			if imgconf, err = png.DecodeConfig(f); err != nil {
				h.Text(w, err.Error(), 200)
				return
			}
		case c_img_jpeg:
			if imgconf, err = jpeg.DecodeConfig(f); err != nil {
				h.Text(w, err.Error(), 200)
				return
			}
		case c_img_gif:
			if imgconf, err = gif.DecodeConfig(f); err != nil {
				h.Text(w, err.Error(), 200)
				return
			}
		default:
			var s string
			if imgconf, s, err = image.DecodeConfig(f); err != nil {
				h.Text(w, err.Error(), 200)
				return
			}
			log.Println(s)
		}
		var img image.Image
		if img, err = png.Decode(f); err != nil {
			h.Text(w, err.Error(), 200)
			return
		}
		h.Text(w, fmt.Sprintln("the img content type is :", imgtype, "img config: ", imgconf.ColorModel, imgconf.Height, imgconf.Width), 200)
		header.Open()
	}
}

