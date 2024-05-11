package utils

import (
	"errors"
	"image"
	"image/png"
	"io"

	ico "github.com/Kodeworks/golang-image-ico"
	"github.com/disintegration/gift"
	"yrh.dev/icns"
)

var ErrUnknownImageType = errors.New("unknown image type")

type ImageType int

const (
	ImageTypeICO ImageType = iota
	ImageTypeICNS
	ImageTypePNG
)

func ConvertPNG(
	input io.Reader,
	output io.Writer,

	imageType ImageType,

	width int,
	height int,
) error {
	src, err := png.Decode(input)
	if err != nil {
		return err
	}

	filter := gift.New(gift.Resize(width, height, gift.LanczosResampling))
	dst := image.NewRGBA(filter.Bounds(src.Bounds()))
	filter.Draw(dst, src) // This does not report any errors

	switch imageType {
	case ImageTypeICO:
		if err := ico.Encode(output, dst); err != nil {
			return err
		}

	case ImageTypeICNS:
		container := icns.NewICNS()

		if err := container.Add(dst); err != nil {
			return err
		}

		if err := icns.Encode(output, container); err != nil {
			return err
		}

	case ImageTypePNG:
		if err := png.Encode(output, dst); err != nil {
			return err
		}

	default:
		return ErrUnknownImageType
	}

	return nil
}
