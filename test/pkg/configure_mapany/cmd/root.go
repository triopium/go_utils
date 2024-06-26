package cmd

import (
	"github.com/triopium/go_utils/pkg/configure"
)

// Build tags set with -ldflags. Cannot set struct fields directly.
var (
	BuildGitTag    string
	BuildGitCommit string
	BuildBuildTime string
)

// VersionInfo Binary version info
var VersionInfo = configure.VersionInfo{
	Version:   "0.0.1",
	GitTag:    BuildGitTag,
	GitCommit: BuildGitCommit,
}

var commandRootConfig = configure.CommanderRoot

func CommandRootRun() {
	commandRootConfig.VersionInfoAdd(&VersionInfo)
	commandRootConfig.AddSub("dummy", RunCommandDummy)
	commandRootConfig.Init()
	commandRootConfig.RunRoot()
}
