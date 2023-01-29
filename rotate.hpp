// +build ignore

#include <Halide.h>
using namespace Halide;

std::tuple<Func, std::vector<Argument>> export_rotate90();
std::tuple<Func, std::vector<Argument>> export_rotate180();
std::tuple<Func, std::vector<Argument>> export_rotate270();
