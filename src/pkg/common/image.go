package common

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
)

type ImageFmt string

const (
	JPEG ImageFmt = "jpeg"
	PNG  ImageFmt = "png"
	GIF  ImageFmt = "gif"
)

func ImageFormat(b []byte) (ImageFmt, error) {
	_, format, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	return ImageFmt(format), nil
}
