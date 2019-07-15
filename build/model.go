package main

import (
	"fmt"
	"path/filepath"
	"runtime"
)

const imagePrefix = "echocat/site24x7_exporter"

var (
	dv = dockerVariant{
		dockerFile: "Dockerfile",
		main:       true,
	}

	currentTarget = target{os: runtime.GOOS, arch: runtime.GOARCH}
	linuxAmd64    = target{os: "linux", arch: "amd64"}
	targets       = []target{
		{os: "darwin", arch: "amd64"},
		{os: "darwin", arch: "386"},
		linuxAmd64,
		{os: "linux", arch: "386"},
		{os: "windows", arch: "amd64"},
		{os: "windows", arch: "386"},
	}
)

type dockerVariant struct {
	dockerFile string
	main       bool
}

func (instance dockerVariant) baseImageName() string {
	return imagePrefix + ":latest"
}

func (instance dockerVariant) imageName(branch string) string {
	result := imagePrefix + ":" + branch
	if branch == "" {
		result += "latest"
	}
	return result
}

type target struct {
	os   string
	arch string
}

func (instance target) outputName() string {
	return filepath.Join("dist", fmt.Sprintf("site24x7_exporter-%s-%s%s", instance.os, instance.arch, instance.ext()))
}

func (instance target) ext() string {
	if instance.os == "windows" {
		return ".exe"
	}
	return ""
}
