package common

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"io"
)

func GetWidthHeightFromImage(r io.Reader) (int, int, string) {
	config, imageType, err := image.DecodeConfig(r)
	Throw_panic(err)

	return config.Width, config.Height, imageType
}
