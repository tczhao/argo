package argo

import (
	"fmt"
	"runtime"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
)

// Version information set by link flags during build. We fall back to these sane
// default values when we build outside the Makefile context (e.g. go build or go test).
var (
	version      = "v0.0.0"               // value from VERSION file
	buildDate    = "1970-01-01T00:00:00Z" // output from `date -u +'%Y-%m-%dT%H:%M:%SZ'`
	gitCommit    = ""                     // output from `git rev-parse HEAD`
	gitTag       = ""                     // output from `git describe --exact-match --tags HEAD` (if clean tree state)
	gitTreeState = ""                     // determined from `git status --porcelain`. either 'clean' or 'dirty'
)

// ImageTag return the image tag.
// GetVersion().Version adulterates the version making it useless as the image tag.
func ImageTag() string {
	if version != "untagged" {
		return version
	}
	return "latest"
}

// GetVersion returns the version information
func GetVersion() wfv1.Version {
	var versionStr = "v3.5.14-atlan-1.2.8"
	return wfv1.Version{
		Version:      versionStr,
		BuildDate:    buildDate,
		GitCommit:    gitCommit,
		GitTag:       gitTag,
		GitTreeState: gitTreeState,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
