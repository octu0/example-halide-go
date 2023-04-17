package example

//go:generate go run ./cmd/download/halide.go

/*
#cgo CFLAGS: -I${SRCDIR}/include
#cgo LDFLAGS: -L${SRCDIR}/lib -ldl -lm
#cgo darwin LDFLAGS: -lrotate90_darwin
#cgo darwin LDFLAGS: -lrotate180_darwin
#cgo darwin LDFLAGS: -lrotate270_darwin
#cgo linux LDFLAGS: -lrotate90_linux
#cgo linux LDFLAGS: -lrotate180_linux
#cgo linux LDFLAGS: -lrotate270_linux

#include "rotate90.h"
#include "rotate180.h"
#include "rotate270.h"
*/
import "C"

import (
	"fmt"
	"image"

	_ "github.com/benesch/cgosymbolizer"
)

//go:generate go run ./cmd/compile/object.go f rotate90 rotate.cpp
func Rotate90(in *image.RGBA) (*image.RGBA, error) {
	width, height := in.Rect.Dx(), in.Rect.Dy()
	out := image.NewRGBA(image.Rect(0, 0, height, width))
	outBuf, err := HalideBufferRGBA(out.Pix, height, width)
	if err != nil {
		return nil, err
	}
	defer HalideFreeBuffer(outBuf)

	inBuf, err := HalideBufferRGBA(in.Pix, width, height)
	if err != nil {
		return nil, err
	}
	defer HalideFreeBuffer(inBuf)

	ret := C.rotate90(
		inBuf,
		C.int(width),
		C.int(height),
		outBuf,
	)
	if ret != C.int(0) {
		return nil, fmt.Errorf("failed to rotate90")
	}
	return out, nil
}

//go:generate go run ./cmd/compile/object.go f rotate180 rotate.cpp
func Rotate180(in *image.RGBA) (*image.RGBA, error) {
	width, height := in.Rect.Dx(), in.Rect.Dy()
	out := image.NewRGBA(image.Rect(0, 0, width, height))
	outBuf, err := HalideBufferRGBA(out.Pix, width, height)
	if err != nil {
		return nil, err
	}
	defer HalideFreeBuffer(outBuf)

	inBuf, err := HalideBufferRGBA(in.Pix, width, height)
	if err != nil {
		return nil, err
	}
	defer HalideFreeBuffer(inBuf)

	ret := C.rotate180(
		inBuf,
		C.int(width),
		C.int(height),
		outBuf,
	)
	if ret != C.int(0) {
		return nil, fmt.Errorf("failed to rotate180")
	}
	return out, nil
}

//go:generate go run ./cmd/compile/object.go f rotate270 rotate.cpp
func Rotate270(in *image.RGBA) (*image.RGBA, error) {
	width, height := in.Rect.Dx(), in.Rect.Dy()
	out := image.NewRGBA(image.Rect(0, 0, height, width))
	outBuf, err := HalideBufferRGBA(out.Pix, height, width)
	if err != nil {
		return nil, err
	}
	defer HalideFreeBuffer(outBuf)

	inBuf, err := HalideBufferRGBA(in.Pix, width, height)
	if err != nil {
		return nil, err
	}
	defer HalideFreeBuffer(inBuf)

	ret := C.rotate270(
		inBuf,
		C.int(width),
		C.int(height),
		outBuf,
	)
	if ret != C.int(0) {
		return nil, fmt.Errorf("failed to rotate90")
	}
	return out, nil
}
