package gospace

import (
	"fmt"
	"path/filepath"
)

var (
	// environment variable containing alternative/gospace specific
	// lookup directories
	SPACES_ENV = "GOSPACES"
	// environment variable containing lookup directories
	CDPATH_ENV = "CDPATH"
)

// resolve the directory against each entry of CDPATH and GOSPACES.
// if the value is already an absolute path and exists in the
// filesystem, it is returned without any further lookups.
func ResolveGospace(dir string) (string, error) {
	if abs, err := filepath.Abs(dir); nil == err {
		if DirExists(abs) {
			D("gospace", dir, "resolves to current working directory")
			return abs, nil
		}
	}

	if abs, ok := SearchPathEnvironment(CDPATH_ENV, dir); ok {
		D("gospace", dir, "was found in", CDPATH_ENV)
		return abs, nil
	} else if abs, ok := SearchPathEnvironment(SPACES_ENV, dir); ok {
		D("gospace", dir, "was found in", SHELL_ENV)
		return abs, nil
	}

	return "", fmt.Errorf("No such directory '%s'", dir)
}
