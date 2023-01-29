// +build ignore

#include <Halide.h>
using namespace Halide;

Func grayscale(Func in, Expr width, Expr height);
Func contrast(Func in, Expr width, Expr height, Expr factor);
Pipeline split(Func in, Expr width, Expr height);

Func export_grayscale();
Func export_contrast();
Pipeline export_split();
