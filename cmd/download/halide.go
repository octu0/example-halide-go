package main

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"os"
	"runtime"
)

func main() {
	switch runtime.GOOS {
	case "darwin":
		downloadDarwin()
	case "linux":
		downloadLinux()
	default:
		panic("does not generate OS type: " + runtime.GOOS)
	}
}

func exists(p string) bool {
	s, err := os.Stat(p)
	if err != nil {
		return false
	}
	if s.IsDir() != true {
		return false
	}
	return true
}

func downloadDarwin() {
	if exists("Halide-14.0.0-x86-64-osx") {
		return
	}

	mustDownload("https://github.com/halide/Halide/releases/download/v14.0.0/Halide-14.0.0-x86-64-osx-6b9ed2afd1d6d0badf04986602c943e287d44e46.tar.gz")
}

func downloadLinux() {
	if exists("Halide-14.0.0-x86-64-linux") {
		return
	}

	mustDownload("https://github.com/halide/Halide/releases/download/v14.0.0/Halide-14.0.0-x86-64-linux-6b9ed2afd1d6d0badf04986602c943e287d44e46.tar.gz")
}

func mustDownload(url string) {
	println("download Halide-14.0.0...")
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	println("extract tar gz...")
	if err := tarxzf(resp.Body); err != nil {
		panic(err)
	}

	println("complete")
}

func tarxzf(r io.Reader) error {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gz.Close()

	t := tar.NewReader(gz)
	for {
		h, err := t.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		if err := tarf(h, t); err != nil {
			return err
		}
	}
	return nil
}

func tarf(h *tar.Header, t *tar.Reader) error {
	switch h.Typeflag {
	case tar.TypeDir:
		if err := os.MkdirAll(h.Name, os.FileMode(h.Mode)); err != nil {
			return err
		}
	case tar.TypeReg:
		f, err := os.Create(h.Name)
		if err != nil {
			return err
		}
		defer f.Close()

		if err := f.Chmod(os.FileMode(h.Mode)); err != nil {
			return err
		}
		if _, err := io.Copy(f, t); err != nil {
			return err
		}
	}
	return nil
}
