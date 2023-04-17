# `example-halide-go`

`example-halide-go` is an example implementation of Halide with `go generate` as a helper.

### How to use

C++ code is written as below:

```cpp
Func rotate90(Func input, Expr width, Expr height) {
  Var x("x"), y("y"), ch("ch");
  Expr w = width - 1;
  Expr h = height - 1;

  Region src_bounds = {{0, w},{0, h},{0, 4}};
  Func in = BoundaryConditions::constant_exterior(input, 0, src_bounds);

  Func r = Func("rotate90");
  r(x, y, ch) = in(y, h - x, ch);
  return r;
}
```

and, Set the needed parameters for the input and output of the function and write them to be returned by std::tuple.

```cpp
std::tuple<Func, std::vector<Argument>> export_rotate90() {
  ImageParam src(UInt(8), 3, "src");
  // input data format
  src.dim(0).set_stride(4);
  src.dim(2).set_stride(1);
  src.dim(2).set_bounds(0, 4);

  Param<int32_t> width{"width", 1920};
  Param<int32_t> height{"height", 1080};

  Func fn = rotate90(src.in(), width, height);

  // output data format
  OutputImageParam out = fn.output_buffer();
  out.dim(0).set_stride(4);
  out.dim(2).set_stride(1);
  out.dim(2).set_bounds(0, 4);

  std::vector<Argument> args = {src, width, height};
  std::tuple<Func, std::vector<Argument>> tuple = std::make_tuple(fn, args);
  return tuple;
}
```

with header(rotate.hpp)

```cpp
#include <Halide.h>
using namespace Halide;

std::tuple<Func, std::vector<Argument>> export_rotate90();
```

Go code is to call it as below:

```go
package example
//go:generate go run ./cmd/download/halide.go

/*
#cgo CFLAGS: -I${SRCDIR}/include
#cgo LDFLAGS: -L${SRCDIR}/lib -ldl -lm -lHalide
#cgo darwin LDFLAGS: -lruntime_darwin
#cgo darwin LDFLAGS: -lrotate90_darwin
#cgo linux LDFLAGS: -lruntime_linux
#cgo linux LDFLAGS: -lrotate90_linux

#include "rotate90.h"
*/
import "C"

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
```
