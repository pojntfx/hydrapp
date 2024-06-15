package utils

import (
	"errors"
	"image"
	"image/png"
	"io"

	ico "github.com/Kodeworks/golang-image-ico"
	"golang.org/x/image/draw"
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

	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil) // This does not report any errors

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
