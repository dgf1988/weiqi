package weiqi

import (
	"fmt"
	"time"
)

const (
	c_ImgBasePath = "img/"
)

/*
CREATE TABLE `img` (
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`title` VARCHAR(100) NOT NULL,
	`name` VARCHAR(100) NOT NULL,
	`md5` CHAR(32) NOT NULL,
	`height` INT(11) NOT NULL,
	`width` INT(11) NOT NULL,
	`type` INT(11) NOT NULL,
	`status` INT(11) NOT NULL DEFAULT '0',
	`upload` DATETIME NOT NULL,
	PRIMARY KEY (`id`)
)
COLLATE='utf8_general_ci'
ENGINE=InnoDB
;
*/

const (
	c_img_png = iota + 1
	c_img_gif
	c_img_bmp
	c_img_jpg
	c_img_jpeg
)

func parseImgType(content_type string) int {
	switch content_type {
	case "image/jpeg":
		return c_img_jpeg
	case "image/gif":
		return c_img_gif
	case "image/bmp":
		return c_img_bmp
	case "image/png":
		return c_img_png
	}
	return 0
}

func formatImgType(img_type int) string {
	switch img_type {
	case c_img_gif:
		return ".gif"
	case c_img_jpeg:
		return ".jpeg"
	case c_img_jpg:
		return ".jpg"
	case c_img_png:
		return ".png"
	case c_img_bmp:
		return ".bmp"
	}
	return ""
}

type Img struct {
	Id     int64
	Title  string
	Name   string
	Md5    string
	Height int64
	Width  int64
	Type   int64
	Status int64
	Upload time.Time
}

func (i Img) GetSuffix() string {
	return formatImgType(int(i.Type))
}

func (i Img) GetPath() string {
	return fmt.Sprintf("%s%s%d/%d/", config.UploadPath, c_ImgBasePath, i.Upload.Year(), i.Upload.Month())
}
func (i Img) GetFilename() string {
	return i.Md5
}

func (i Img) Alt() string {
	return i.Title
}

func (i Img) Src() string {
	return fmt.Sprintf("/%s%d/%d/%s%s", c_ImgBasePath, i.Upload.Year(), i.Upload.Month(), i.GetFilename(), i.GetSuffix())
}

func (i Img) GetFullname() string { return fmt.Sprint(i.GetPath(), i.GetFilename(), i.GetSuffix()) }
