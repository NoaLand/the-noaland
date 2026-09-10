package image

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"strings"
	"testing"
)

func TestRenderKittyImagePayload(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	img.SetNRGBA(0, 0, color.NRGBA{R: 255, A: 255})
	output := Render(img, 14, 7)
	header := "\x1b_Ga=T,f=100,c=14,r=7,q=2,m=0;"
	if !strings.HasPrefix(output, header) || !strings.HasSuffix(output, "\x1b\\") {
		t.Fatal("invalid Kitty dimensions or framing")
	}
	payload := strings.TrimSuffix(strings.TrimPrefix(output, header), "\x1b\\")
	data, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if decoded.Bounds() != img.Bounds() || color.NRGBAModel.Convert(decoded.At(0, 0)) != img.NRGBAAt(0, 0) {
		t.Fatal("image dimensions or pixels changed")
	}
}
