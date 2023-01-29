package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/octu0/example-halide-go"

	_ "embed"
)

var (
	//go:embed src.png
	data []byte
)

func main() {
	src, err := pngToRGBA(data)
	if err != nil {
		panic(err)
	}

	if out, err := example.Rotate90(src); err != nil {
		panic(err)
	} else {
		file, err := saveImage(out)
		if err != nil {
			panic(err)
		}
		println("rotate90:", file)
	}

	if out, err := example.Rotate180(src); err != nil {
		panic(err)
	} else {
		file, err := saveImage(out)
		if err != nil {
			panic(err)
		}
		println("rotate180:", file)
	}

	if out, err := example.Rotate270(src); err != nil {
		panic(err)
	} else {
		file, err := saveImage(out)
		if err != nil {
			panic(err)
		}
		println("rotate270:", file)
	}
}

func pngToRGBA(data []byte) (*image.RGBA, error) {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	if i, ok := img.(*image.RGBA); ok {
		return i, nil
	}

	b := img.Bounds()
	rgba := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y += 1 {
		for x := b.Min.X; x < b.Max.X; x += 1 {
			c := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
			rgba.Set(x, y, c)
		}
	}
	return rgba, nil
}

func saveImage(img *image.RGBA) (string, error) {
	out, err := os.CreateTemp("/tmp", "out*.png")
	if err != nil {
		return "", err
	}
	defer out.Close()

	if err := png.Encode(out, img); err != nil {
		return "", err
	}
	return out.Name(), nil
}
