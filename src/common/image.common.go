package common

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
)

func GetImageFormat(b []byte) (string, error) {
	_, format, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	return format, nil
}
