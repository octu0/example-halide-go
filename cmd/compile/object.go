package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

const generateRuntimeMainTmpl string = `
#include <Halide.h>
using namespace Halide;
int main() {
  Var x("x");
  Func fn = Func("noop");
  fn(x) = 0;
  std::vector<Argument> args;

  std::vector<Target::Feature> features;
  features.push_back(Target::AVX);
  features.push_back(Target::AVX2);
  features.push_back(Target::FMA);
  features.push_back(Target::F16C);
  features.push_back(Target::SSE41);
  {
    Target target;
    target.os = Target::OSX;
    target.arch = Target::X86;
    target.bits = 64;
    target.set_features(features);
    fn.compile_to_object("{{ .FileNameDarwin }}", args, "{{ .Name }}", target);
  }
  {
    Target target;
    target.os = Target::Linux;
    target.arch = Target::X86;
    target.bits = 64;
    target.set_features(features);
    fn.compile_to_object("{{ .FileNameLinux }}", args, "{{ .Name }}", target);
  }
  fn.compile_to_header("{{ .HeaderName }}", args, "{{ .Name }}");
  return 0;
}
`

const generateGenRunMainTmpl string = `
#include <Halide.h>
#include "{{ .HppFileName }}"
using namespace Halide;
int main() {
  std::vector<Target::Feature> features;
  features.push_back(Target::AVX);
  features.push_back(Target::AVX2);
  features.push_back(Target::FMA);
  features.push_back(Target::F16C);
  features.push_back(Target::SSE41);
  features.push_back(Target::Feature::NoRuntime);

  std::tuple<{{ .ExportType }}, std::vector<Argument>> tuple = export_{{ .Name }}();
  {{ .ExportType }} fn = std::get<0>(tuple);
  std::vector<Argument> args =std::get<1>(tuple);
  {
    Target target;
    target.os = Target::OSX;
    target.arch = Target::X86;
    target.bits = 64;
    target.set_features(features);
    fn.compile_to_object("{{ .FileNameDarwin }}", args, "{{ .Name }}", target);
  }
  {
    Target target;
    target.os = Target::Linux;
    target.arch = Target::X86;
    target.bits = 64;
    target.set_features(features);
    fn.compile_to_object("{{ .FileNameLinux }}", args, "{{ .Name }}", target);
  }
  fn.compile_to_header("{{ .HeaderName }}", args, "{{ .Name }}");
  return 0;
}
`

type GenRun struct {
	FileNameDarwin    string
	FileNameLinux     string
	AsmNameDarwin     string
	AsmNameLinux      string
	LLVMAsmNameDarwin string
	LLVMAsmNameLinux  string
	LLVMBcNameDarwin  string
	LLVMBcNameLinux   string
	HeaderName        string
	Name              string
	HppFileName       string
	ExecFileName      string
	MainTemplate      string
	ExportType        string
}

func main() {
	halidePath := ""
	switch runtime.GOOS {
	case "darwin":
		halidePath = "Halide-14.0.0-x86-64-osx"
	case "linux":
		halidePath = "Halide-14.0.0-x86-64-linux"
	default:
		panic("not support os: " + runtime.GOOS)
	}

	exportTypeName := os.Args[1]
	funcName := os.Args[2]
	targetFilePath := os.Args[3]
	targetFileBase := filepath.Base(targetFilePath)
	targetFileExt := filepath.Ext(targetFileBase)
	baseName := targetFileBase[0:strings.LastIndex(targetFileBase, targetFileExt)]

	exportType := ""
	switch exportTypeName {
	case "f", "func":
		exportType = "Func"
	case "p", "pipeline":
		exportType = "Pipeline"
	}
	if exportType == "" {
		panic("not support extern type: " + exportTypeName)
	}

	targets := make([]GenRun, 2)
	targets[0] = GenRun{
		FileNameDarwin:   fmt.Sprintf("lib/lib%s_darwin.dylib", "runtime"),
		FileNameLinux:    fmt.Sprintf("lib/lib%s_linux.o", "runtime"),
		AsmNameDarwin:    fmt.Sprintf("lib/lib%s_darwin.s", "runtime"),
		AsmNameLinux:     fmt.Sprintf("lib/lib%s_linux.s", "runtime"),
		LLVMBcNameDarwin: fmt.Sprintf("lib/lib%s_darwin.bc", "runtime"),
		LLVMBcNameLinux:  fmt.Sprintf("lib/lib%s_linux.bc", "runtime"),
		HeaderName:       fmt.Sprintf("include/%s.h", "runtime"),
		Name:             "runtime",
		HppFileName:      "",
		ExecFileName:     fmt.Sprintf("gen/%s.out", "runtime"),
		MainTemplate:     generateRuntimeMainTmpl,
		ExportType:       "",
	}
	targets[1] = GenRun{
		FileNameDarwin:   fmt.Sprintf("lib/lib%s_darwin.dylib", funcName),
		FileNameLinux:    fmt.Sprintf("lib/lib%s_linux.o", funcName),
		AsmNameDarwin:    fmt.Sprintf("lib/lib%s_darwin.s", funcName),
		AsmNameLinux:     fmt.Sprintf("lib/lib%s_linux.s", funcName),
		LLVMBcNameDarwin: fmt.Sprintf("lib/lib%s_darwin.bc", funcName),
		LLVMBcNameLinux:  fmt.Sprintf("lib/lib%s_linux.bc", funcName),
		HeaderName:       fmt.Sprintf("include/%s.h", funcName),
		Name:             funcName,
		HppFileName:      fmt.Sprintf("%s.hpp", baseName),
		ExecFileName:     fmt.Sprintf("gen/%s.out", funcName),
		MainTemplate:     generateGenRunMainTmpl,
		ExportType:       exportType,
	}

	libpng := exec.Command("libpng-config", "--cflags", "--ldflags")
	libpngCfg, err := libpng.Output()
	if err != nil {
		panic(err)
	}
	libpngFlags := strings.TrimSpace(string(libpngCfg))
	libpngFlags = strings.ReplaceAll(libpngFlags, "\n", " ")

	buf := bytes.NewBuffer(nil)
	files := make([]*os.File, 0, len(targets))
	defer func() {
		for _, f := range files {
			os.Remove(f.Name())
		}
	}()
	for _, t := range targets {
		if _, err := os.Stat(t.ExecFileName); err == nil {
			continue // file exists
		}
		buf.Reset()

		println("compiling...", t.Name)
		tpl, err := template.New(t.Name).Parse(t.MainTemplate)
		if err != nil {
			panic(err)
		}
		if err := tpl.Execute(buf, t); err != nil {
			panic(err)
		}
		mainC, err := os.CreateTemp("", "main*.cpp")
		if err != nil {
			panic(err)
		}
		files = append(files, mainC)
		if _, err := mainC.Write(buf.Bytes()); err != nil {
			panic(err)
		}

		genArgs := []string{
			"clang++",
			"-g",
			"-I.",
			"-I" + halidePath + "/include",
			"-I" + halidePath + "/share/Halide/tools",
			"-L" + halidePath + "/lib",
			libpngFlags,
			"-L/usr/local/opt/jpeg/lib",
			"-I/usr/local/opt/jpeg/include",
			"-ljpeg",
			"-lHalide",
			"-lpthread",
			"-ldl",
			"-lz",
			"-std=c++17",
			"-o", t.ExecFileName,
		}
		if t.Name == "runtime" {
			genArgs = append(genArgs, mainC.Name())
		} else {
			genArgs = append(genArgs, targetFilePath)
			genArgs = append(genArgs, mainC.Name())
		}

		println("compile...", t.Name, "cmd:", strings.Join(genArgs, " "))

		cmd := exec.Command("sh", "-c", strings.Join(genArgs, " "))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			println("error", err.Error())
			continue
		}

		println("generate...", t.Name)
		gen := exec.Command("sh", "-c", t.ExecFileName)
		gen.Stdout = os.Stdout
		gen.Stderr = os.Stderr
		if err := gen.Run(); err != nil {
			println("error", err.Error())
			continue
		}
		println("done...", t.Name)
	}
}
