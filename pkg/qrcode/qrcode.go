package qrcode

import (
	"NULL/consul-gin/pkg/file"
	"NULL/consul-gin/pkg/setting"
	"NULL/consul-gin/pkg/util"
	"errors"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/nfnt/resize"
	"github.com/skip2/go-qrcode"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"time"
)

type QrCode struct {
	URL    string
	Width  int
	Height int
	Ext    string
	Level  qr.ErrorCorrectionLevel
	Mode   qr.Encoding
}

const (
	logo_file = "bg.jpg"
	logo_w    = 54
	logo_h    = 54
	qr_size   = 256
	qr_level  = qrcode.High
	EXT_JPG   = ".jpg"
)

// NewQrCode initialize instance
func NewQrCode(url string, width, height int, level qr.ErrorCorrectionLevel, mode qr.Encoding) *QrCode {
	return &QrCode{
		URL:    url,
		Width:  width,
		Height: height,
		Level:  level,
		Mode:   mode,
		Ext:    EXT_JPG,
	}
}

// GetQrCodePath get save path
func GetQrCodePath() string {
	return setting.AppSetting.QrCodeSavePath
}

// GetQrCodeFullPath get full save path
func GetQrCodeFullPath() string {
	return setting.AppSetting.RuntimeRootPath + setting.AppSetting.QrCodeSavePath
}

// GetQrCodeFullUrl get the full access path
func GetQrCodeFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetQrCodePath() + name
}

// GetQrCodeFileName get qr file name
func GetQrCodeFileName(value string) string {
	return util.EncodeMD5(value)
}

// GetQrCodeExt get qr file ext
func (q *QrCode) GetQrCodeExt() string {
	return q.Ext
}

// Encode generate QR code
func (q *QrCode) Encode(path string) (string, string, error) {
	//name := GetQrCodeFileName(q.URL) + q.GetQrCodeExt()
	date := time.Now().Format("20060102")
	name := q.URL + date + q.GetQrCodeExt()
	src := path + name
	if file.CheckNotExist(src) == true {
		code, err := qr.Encode(q.URL, q.Level, q.Mode)
		if err != nil {
			return "", "", err
		}

		code, err = barcode.Scale(code, q.Width, q.Height)
		if err != nil {
			return "", "", err
		}

		f, err := file.MustOpen(name, path)
		if err != nil {
			return "", "", err
		}
		defer f.Close()

		err = jpeg.Encode(f, code, nil)
		if err != nil {
			return "", "", err
		}
	}

	return name, path, nil
}

// Encode generate QR code with logo
func GenerateQrWithLogo(uri, path string) (string, string, error) {
	var (
		err error
		q   *qrcode.QRCode
	)
	date := time.Now().Format("20060102150405")
	name := uri + date + EXT_JPG
	src := path + name
	if file.CheckNotExist(src) == true {
		// 先创建一个二维码
		q, err = qrcode.New(uri, qr_level)
		if err != nil {
			err = errors.New("can not create a qrcode")
			return "", "", err
		}
		png := q.Image(qr_size)
		bounds := png.Bounds()
		// 通过二维码创建一个空画布
		newImg := image.NewRGBA(bounds)
		// 读取logo文件
		logo_path := setting.AppSetting.RuntimeRootPath + logo_file
		filelogo, err := os.Open(logo_path)
		if err != nil {
			return "", "", err
		}
		defer filelogo.Close()
		// 将file转换成image
		logo, _, err := image.Decode(filelogo)
		if err != nil {
			return "", "", err
		}
		// 按比例缩放logo
		logo = resize.Resize(logo_w, 0, logo, resize.Lanczos3)
		// 在画布上分别画上二维码，缩略后的logo
		draw.Draw(newImg, newImg.Bounds(), png, png.Bounds().Min, draw.Over)
		draw.Draw(newImg, image.Rect((qr_size/2)-(logo_w/2), (qr_size/2)-(logo_h/2), (qr_size/2)+(logo_w/2), (qr_size/2)+(logo_h/2)), logo, logo.Bounds().Min, draw.Over)
		f, err := file.MustOpen(name, path)
		if err != nil {
			return "", "", err
		}
		defer f.Close()
		err = jpeg.Encode(f, newImg, nil)
		if err != nil {
			return "", "", err
		}
	}
	return name, path, nil
}
