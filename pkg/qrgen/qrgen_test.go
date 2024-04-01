package qrgen_test

import (
	"bytes"
	"encoding/base64"
	"image"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gosom/toolkit/pkg/qrgen"
)

func TestQRGen_GenerateQRCode(t *testing.T) {
	gen := qrgen.NewQRGen()

	buf := new(bytes.Buffer)
	err := gen.GenerateQRCode(buf, "test", 0)
	require.NoError(t, err)

	require.NotEmpty(t, buf.Bytes())
}

func TestQRGen_GenerateQRCodeAsBase64(t *testing.T) {
	gen := qrgen.NewQRGen()

	base64Str, err := gen.GenerateQRCodeAsBase64("https://haviohealth.com", 0)
	require.NoError(t, err)
	require.NotEmpty(t, base64Str)

	if idx := strings.Index(base64Str, "base64,"); idx != -1 {
		base64Str = base64Str[idx+7:]
	}

	imgData, err := base64.StdEncoding.DecodeString(base64Str)
	require.NoError(t, err)

	_, _, err = image.Decode(strings.NewReader(string(imgData)))
	require.NoError(t, err)
}
