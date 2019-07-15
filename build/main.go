package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
)

var (
	branch                    = "snapshot"
	commit                    = "unknown"
	withDocker                = true
	plainLatestVersionPattern string
	latestVersionPattern      *regexp.Regexp
)

func init() {
	flag.StringVar(&branch, "branch", os.Getenv("TRAVIS_BRANCH"), "something like either main, v1.2.3 or snapshot-feature-foo")
	flag.StringVar(&commit, "commit", os.Getenv("TRAVIS_COMMIT"), "something like 463e189796d5e96a7b605ab51985458faf8fd0d4")
	flag.BoolVar(&withDocker, "withDocker", true, "enables docker tests and builds")
	flag.StringVar(&plainLatestVersionPattern, "latestVersionPattern", os.Getenv("LATEST_VERSION_PATTERN"), "everything what matches here will be a latest tag")
}

func main() {
	flag.Parse()
	if branch == "" {
		branch = "development"
	}
	if commit == "" {
		commit = "development"
	}
	if compiled, err := regexp.Compile(plainLatestVersionPattern); err != nil {
		panic(err)
	} else {
		latestVersionPattern = compiled
	}
	fArgs := flag.Args()
	if len(fArgs) != 1 {
		usage()
	}
	switch fArgs[0] {
	case "build":
		build(branch, commit)
	case "test":
		test(branch, commit)
	case "deploy":
		deploy(branch)
	case "build-and-deploy":
		build(branch, commit)
		deploy(branch)
	default:
		usage()
	}
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, "Usage: %s [flags] <command>", os.Args[0])
	os.Exit(1)
}
