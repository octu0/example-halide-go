// +build ignore

#include <Halide.h>
#include "rotate.hpp"

using namespace Halide;

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

Func rotate180(Func input, Expr width, Expr height) {
  Var x("x"), y("y"), ch("ch");
  Expr w = width - 1;
  Expr h = height - 1;

  Region src_bounds = {{0, w},{0, h},{0, 4}};
  Func in = BoundaryConditions::constant_exterior(input, 0, src_bounds);

  Func r = Func("rotate180");
  r(x, y, ch) = in(w - x, h - y, ch);
  return r;
}

Func rotate270(Func input, Expr width, Expr height) {
  Var x("x"), y("y"), ch("ch");
  Expr w = width - 1;
  Expr h = height - 1;

  Region src_bounds = {{0, w},{0, h},{0, 4}};
  Func in = BoundaryConditions::constant_exterior(input, 0, src_bounds);

  Func r = Func("rotate270");
  r(x, y, ch) = in(w - y, x, ch);
  return r;
}

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

std::tuple<Func, std::vector<Argument>> export_rotate180() {
  ImageParam src(UInt(8), 3, "src");
  // input data format
  src.dim(0).set_stride(4);
  src.dim(2).set_stride(1);
  src.dim(2).set_bounds(0, 4);

  Param<int32_t> width{"width", 1920};
  Param<int32_t> height{"height", 1080};

  Func fn = rotate180(src.in(), width, height);

  // output data format
  OutputImageParam out = fn.output_buffer();
  out.dim(0).set_stride(4);
  out.dim(2).set_stride(1);
  out.dim(2).set_bounds(0, 4);

  std::vector<Argument> args = {src, width, height};
  std::tuple<Func, std::vector<Argument>> tuple = std::make_tuple(fn, args);
  return tuple;
}

std::tuple<Func, std::vector<Argument>> export_rotate270() {
  ImageParam src(UInt(8), 3, "src");
  // input data format
  src.dim(0).set_stride(4);
  src.dim(2).set_stride(1);
  src.dim(2).set_bounds(0, 4);

  Param<int32_t> width{"width", 1920};
  Param<int32_t> height{"height", 1080};

  Func fn = rotate270(src.in(), width, height);

  // output data format
  OutputImageParam out = fn.output_buffer();
  out.dim(0).set_stride(4);
  out.dim(2).set_stride(1);
  out.dim(2).set_bounds(0, 4);

  std::vector<Argument> args = {src, width, height};
  std::tuple<Func, std::vector<Argument>> tuple = std::make_tuple(fn, args);
  return tuple;
}
