package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func test(branch, commit string) {
	testGoCode(currentTarget)

	buildBinary(branch, commit, currentTarget, true)
	testBinary(branch, commit, currentTarget)

	if withDocker {
		buildBinary(branch, commit, linuxAmd64, true)
		buildDocker(branch, dv, true)
		testDocker(branch, commit, dv)
		tagDocker(branch, dv, true)
	}
}

func testGoCode(t target) {
	executeTo(func(cmd *exec.Cmd) {
		cmd.Env = append(os.Environ(), "GOOS="+t.os, "GOARCH="+t.arch)
	}, os.Stderr, os.Stdout, "go", "test", "-v", "./...")
}

func testBinary(branch, commit string, t target) {
	testBinaryByExpectingResponse(t, `site24x7_exporter (version: TEST`+branch+`TEST, revision: TEST`+commit+`TEST, build: `, t.outputName(), "-help")
}

func testBinaryByExpectingResponse(t target, expectedPartOfResponse string, args ...string) {
	cmd := append([]string{t.outputName()}, args...)
	response := executeAndRecord(args...)
	if !strings.Contains(response, expectedPartOfResponse) {
		panic(fmt.Sprintf("Command failed [%s]\nResponse should contain: %s\nBut response was: %s",
			quoteAllIfNeeded(cmd...), expectedPartOfResponse, response))
	}
}

func testDocker(branch, commit string, v dockerVariant) {
	testDockerByExpectingResponse(branch, v, `site24x7_exporter (version: TEST`+branch+`TEST, revision: TEST`+commit+`TEST, build: `, "-help")
}

func testDockerByExpectingResponse(branch string, v dockerVariant, expectedPartOfResponse string, command ...string) {
	call := append([]string{"docker", "run", "--rm", v.imageName("TEST" + branch + "TEST")}, command...)
	response := executeAndRecord(call...)
	if !strings.Contains(response, expectedPartOfResponse) {
		panic(fmt.Sprintf("Command failed [%s]\nResponse should contain: %s\nBut response was: %s",
			quoteAllIfNeeded(call...), expectedPartOfResponse, response))
	}
}
