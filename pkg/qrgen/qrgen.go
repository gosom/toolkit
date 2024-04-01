package qrgen

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"sync"

	"github.com/skip2/go-qrcode"

	"toolkit/pkg/errorsext"
	"toolkit/pkg/logger"
)

const (
	// QRCodeSize is the default size of the QR code
	QRCodeSize = 128
)

var (
	qOnce        sync.Once
	defaultQRGen *QRGen
)

func GetQRCode(data string) string {
	qOnce.Do(func() {
		defaultQRGen = NewQRGen()
	})

	ans, err := defaultQRGen.GenerateQRCodeAsBase64(data, 0)
	if err != nil {
		logger.Error(context.Background(), "error generating QR", "error", err)

		return ""
	}

	return ans
}

type QRGen struct {
}

func NewQRGen() *QRGen {
	return &QRGen{}
}

func (q *QRGen) GenerateQRCode(w io.Writer, data string, size int) error {
	if size == 0 {
		size = QRCodeSize
	}

	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return errorsext.WithStack(err)
	}

	err = qr.Write(size, w)
	if err != nil {
		return errorsext.WithStack(err)
	}

	return nil
}

func (q *QRGen) GenerateQRCodeAsBase64(data string, size int) (string, error) {
	if size == 0 {
		size = QRCodeSize
	}

	imgData, err := qrcode.Encode(data, qrcode.Medium, size)
	if err != nil {
		return "", errorsext.WithStack(err)
	}

	imgData = bytes.TrimSpace(imgData)

	base64Str := base64.StdEncoding.EncodeToString(imgData)
	webBase64 := "data:image/png;base64," + base64Str

	return webBase64, nil
}
