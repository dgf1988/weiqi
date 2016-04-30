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
	`height` INT(11) NOT NULL,
	`width` INT(11) NOT NULL,
	`type` INT(11) NOT NULL DEFAULT '0',
	`status` INT(11) NOT NULL DEFAULT '0',
	`upload` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`)
)
ENGINE=InnoDB
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

func parseTypeSuffix(img_type int) string {
	switch img_type {
	case c_img_bmp:
		return ".bmp"
	case c_img_gif:
		return ".gif"
	}
}

type Img struct {
	Id     int64
	Title  string
	Height int64
	Width  int64
	Type   int64
	Status int64
	Upload time.Time
}

func (i Img) GetSuffix() string {
	switch i.Type {
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

func (i Img) GetPath() string {	return fmt.Sprintf("%s/%s/", i.Upload.Year(), i.Upload.Month()) }
func (i Img) GetFilename() string {	return getMd5(fmt.Sprint(i.Title, i.Upload.UnixNano())) }
func (i Img) GetSrc() string { return fmt.Sprint(config.UploadPath, c_ImgBasePath, i.GetPath(), i.GetFilename(), i.GetSuffix()) }
