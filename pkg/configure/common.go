package configure

import "fmt"

// VersionInfo
type VersionInfo struct {
	ProgramName string
	Version     string
	GitTag      string
	GitCommit   string
	BuildTime   string
}

var FlagsUsage = "Usage:\n"

// Usage called when help command invoked
func Usage() {
	fmt.Println(FlagsUsage)
}
