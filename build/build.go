package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func build(branch, commit string) {
	buildBinaries(branch, commit)

	if withDocker {
		buildDocker(branch, dv, false)
		tagDocker(branch, dv, false)
	}
}

func buildBinaries(branch, commit string) {
	for _, t := range targets {
		buildBinary(branch, commit, t, false)
	}
}

func buildBinary(branch, commit string, t target, forTesting bool) {
	ldFlags := buildLdFlagsFor(branch, commit, forTesting)
	outputName := t.outputName()
	must(os.MkdirAll(filepath.Dir(outputName), 0755))
	executeTo(func(cmd *exec.Cmd) {
		cmd.Env = append(os.Environ(), "GOOS="+t.os, "GOARCH="+t.arch)
	}, os.Stderr, os.Stdout, "go", "build", "-ldflags", ldFlags, "-o", outputName, ".")
}

func buildLdFlagsFor(branch, commit string, forTesting bool) string {
	testPrefix := ""
	testSuffix := ""
	if forTesting {
		testPrefix = "TEST"
		testSuffix = "TEST"
	}
	return fmt.Sprintf("-X main.version=%s%s%s", testPrefix, branch, testSuffix) +
		fmt.Sprintf(" -X main.revision=%s%s%s", testPrefix, commit, testSuffix) +
		fmt.Sprintf(" -X main.compiled=%s", startTime.Format("2006-01-02T15:04:05Z"))
}

func buildDocker(branch string, v dockerVariant, forTesting bool) {
	version := branch
	if forTesting {
		version = "TEST" + version + "TEST"
	}
	execute("docker", "build", "-t", v.imageName(version), "-f", v.dockerFile, "--build-arg", "image="+imagePrefix, "--build-arg", "version="+version, ".")
}

func tagDocker(branch string, v dockerVariant, forTesting bool) {
	version := branch
	if forTesting {
		version = "TEST" + version + "TEST"
	}
	executeForVersionParts(version, func(tagSuffix string) {
		tagDockerWith(version, v, v.imageName(tagSuffix))
	})
	if latestVersionPattern != nil && latestVersionPattern.MatchString(version) {
		tagDockerWith(version, v, v.baseImageName())
	}
	if v.main {
		tagDockerWith(version, v, imagePrefix+":"+version)
		executeForVersionParts(version, func(tagSuffix string) {
			tagDockerWith(version, v, imagePrefix+":"+tagSuffix)
		})
		if latestVersionPattern != nil && latestVersionPattern.MatchString(version) {
			tagDockerWith(version, v, imagePrefix+":latest")
		}
	}
}

func tagDockerWith(branch string, v dockerVariant, tag string) {
	execute("docker", "tag", v.imageName(branch), tag)
}
