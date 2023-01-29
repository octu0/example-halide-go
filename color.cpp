// +build ignore

#include <Halide.h>
#include "color.hpp"

using namespace Halide;

const uint8_t GRAY_R = 76, GRAY_G = 152, GRAY_B = 28;

Func grayscale(Func input, Expr width, Expr height) {
  Var x("x"), y("y"), ch("ch");
  Var xo("xo"), xi("xi");
  Var yo("yo"), yi("yi");
  Var ti("ti");

  Region src_bounds = {{0, width},{0, height},{0, 4}};
  Func in = BoundaryConditions::repeat_edge(input, src_bounds);

  Func fn = Func("grayscale");
  Expr r = cast<int16_t>(in(x, y, 0));
  Expr g = cast<int16_t>(in(x, y, 1));
  Expr b = cast<int16_t>(in(x, y, 2));
  Expr a = cast<int16_t>(in(x, y, 3));
  Expr value = ((r * GRAY_R) + (g * GRAY_G) + (b * GRAY_B)) >> 8;
  value = cast<uint8_t>(value);

  fn(x, y, ch) = cast<uint8_t>(255);
  fn(x, y, 0) = value; 
  fn(x, y, 1) = value; 
  fn(x, y, 2) = value; 
  fn(x, y, 3) = cast<uint8_t>(a); 

  // schedule
  fn.update(0).unscheduled();
  fn.update(1).unscheduled();
  fn.update(2).unscheduled();
  fn.update(3).unscheduled();
  fn.compute_at(in, ti)
    .tile(x, y, xo, yo, xi, yi, 32, 32)
    .fuse(xo, yo, ti)
    .parallel(ch)
    .parallel(ti, 8)
    .vectorize(xi, 32);
  return fn;
}

Func contrast(Func input, Expr width, Expr height, Expr factor) {
  Var x("x"), y("y"), ch("ch");

  Region src_bounds = {{0, width},{0, height},{0, 4}};
  Func in = BoundaryConditions::repeat_edge(input, src_bounds);

  Expr e = max(min(1.0f, factor), 0.0f);

  Func contrast = Func("contrast");
  Expr value = in(x, y, ch);
  value = cast<float>(value);
  value = (value / 255.0f) - 0.5f;
  value = (value * e) + 0.5f;
  value = value * 255.0f;

  contrast(x, y, ch) = cast<uint8_t>(value);

  // schedule
  contrast.compute_at(in, x)
    .parallel(ch)
    .vectorize(x, 16);
  return contrast;
}

Pipeline split(Func input, Expr width, Expr height) {
  Var x("x"), y("y"), ch("ch");

  Region src_bounds = {{0, width},{0, height},{0, 4}};
  Func in = BoundaryConditions::repeat_edge(input, src_bounds);

  Func r = Func("split_red");
  Func g = Func("split_green");
  Func b = Func("split_blue");

  r(x, y, ch) = undef(UInt(8));
  r(x, y, 0) = cast<uint8_t>(in(x, y, 0));
  r(x, y, 1) = cast<uint8_t>(0);
  r(x, y, 2) = cast<uint8_t>(0);
  r(x, y, 3) = cast<uint8_t>(255);

  g(x, y, ch) = undef(UInt(8));
  g(x, y, 0) = cast<uint8_t>(0);
  g(x, y, 1) = cast<uint8_t>(in(x, y, 1));
  g(x, y, 2) = cast<uint8_t>(0);
  g(x, y, 3) = cast<uint8_t>(255);

  b(x, y, ch) = undef(UInt(8));
  g(x, y, 0) = cast<uint8_t>(0);
  g(x, y, 1) = cast<uint8_t>(0);
  g(x, y, 2) = cast<uint8_t>(in(x, y, 2));
  g(x, y, 3) = cast<uint8_t>(255);

  return Pipeline({r, g, b});
}

std::tuple<Func, std::vector<Argument>> export_grayscale() {
  ImageParam src(UInt(8), 3, "src");
  // input data format
  src.dim(0).set_stride(4);
  src.dim(2).set_stride(1);
  src.dim(2).set_bounds(0, 4);

  Param<int32_t> width{"width", 1920};
  Param<int32_t> height{"height", 1080};

  Func fn = grayscale(src.in(), width, height);

  // output data format
  OutputImageParam out = fn.output_buffer();
  out.dim(0).set_stride(4);
  out.dim(2).set_stride(1);
  out.dim(2).set_bounds(0, 4);

  std::vector<Argument> args = {src, width, height};
  std::tuple<Func, std::vector<Argument>> tuple = std::make_tuple(fn, args);
  return tuple;
}

std::tuple<Func, std::vector<Argument>> export_contrast() {
  ImageParam src(UInt(8), 3, "src");
  // input data format
  src.dim(0).set_stride(4);
  src.dim(2).set_stride(1);
  src.dim(2).set_bounds(0, 4);

  Param<int32_t> width{"width", 1920};
  Param<int32_t> height{"height", 1080};
  Param<float> factor{"factor", 0.525f};

  Func fn = contrast(src.in(), width, height, factor);

  // output data format
  OutputImageParam out = fn.output_buffer();
  out.dim(0).set_stride(4);
  out.dim(2).set_stride(1);
  out.dim(2).set_bounds(0, 4);

  std::vector<Argument> args = {src, width, height, factor};
  std::tuple<Func, std::vector<Argument>> tuple = std::make_tuple(fn, args);
  return tuple;
}

std::tuple<Pipeline, std::vector<Argument>> export_split() {
  ImageParam src(UInt(8), 3, "src");
  // input data format
  src.dim(0).set_stride(4);
  src.dim(2).set_stride(1);
  src.dim(2).set_bounds(0, 4);

  Param<int32_t> width{"width", 1920};
  Param<int32_t> height{"height", 1080};

  Pipeline fn = split(src.in(), width, height);

  std::vector<Argument> args = {src, width, height, factor};
  std::tuple<Func, std::vector<Argument>> tuple = std::make_tuple(fn, args);
  return tuple;
}
