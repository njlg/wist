// Package info provides info about this app.
package info

// GetVersion returns the version of the executable.
func GetVersion() string {
	return "local-dev"
}

// GetSHA returns the git SHA that was used when the binary was built.
func GetSHA() string {
	return "master"
}
