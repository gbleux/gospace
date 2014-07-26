package gospace

import (
	"os"
	"path/filepath"
	"strings"
)

// searches for an existing path in the envionment. the first parameter
// is expected to be an environment variable containing path-list
// separated entries. each enty itself is converted to an absolute
// path first and then joined with the relative parameter. if the path
// resolves to an actual node on the filesystem, is is returned. if
// none of the lookups yielded a result, an empty string is returned.
//
// if the second parameter is already absolute it is returned if it
// refers to an existing node, otherwise it is treated relative to
// each path entry of the environment variable.
func SearchPathEnvironment(env string, rel string) (string, bool) {
	path := os.Getenv(env)
	fragments := strings.Split(path, string(os.PathListSeparator))

	// early exit for existing absolute path
	if filepath.IsAbs(rel) && PathExists(rel) {
		return rel, true
	}

	T("searching for", rel, "in", env)

	for _, fragment := range fragments {
		node := filepath.Join(fragment, rel)

		// any error or negative result will continue the loop
		if abs, err := filepath.Abs(node); nil == err {
			if PathExists(abs) {
				return abs, true
			}
		}
	}

	return "", false
}

// check if the provided path resembles an actual node
// in the filesystem. the return value can also indicate a lack of
// access privileges or other problems.
func PathExists(path string) bool {
	_, err := os.Stat(path)

	return nil == err
}

// similar to PathExists, but also check if the filesystem node is
// a directory
func DirExists(path string) bool {
	if info, _ := os.Stat(path); nil != info {
		return info.IsDir()
	}

	return false
}
