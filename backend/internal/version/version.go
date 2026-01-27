package version

// These variables are set at build time using -ldflags
var (
	// GitCommit is the git commit hash
	GitCommit = "development"
	// GitBranch is the git branch name
	GitBranch = "unknown"
	// BuildTime is the build timestamp
	BuildTime = "unknown"
)

// Short returns the first 7 characters of the git commit hash
func Short() string {
	if len(GitCommit) >= 7 {
		return GitCommit[:7]
	}
	return GitCommit
}

// Info returns version information as a map
func Info() map[string]string {
	return map[string]string{
		"commit":     GitCommit,
		"short":      Short(),
		"branch":     GitBranch,
		"build_time": BuildTime,
	}
}
