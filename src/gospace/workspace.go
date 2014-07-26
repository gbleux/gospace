package gospace

import (
	"os"
	"path"
	"strings"
)

const (
	// sub-directory to include in PATH and GOBIN
	BIN_DIR string = "bin"
	// environment variable pointing to a golang installation
	SDK_ENV string = "GOHOME"
	// workspace environment variable for binaries
	PKG_ENV string = "GOBIN"
	// environment variable containing lookup directories
	OS_ENV string = "PATH"
	// environment variable containing lookup directories
	WS_ENV string = "GOPATH"
)

var (
	// default workspace directory
	WS_DEFAULT string
)

// GO workspace paths
type Workspace struct {
	// working directory
	Root string
	// include directories
	GoPath []string
	// OS PATH directories
	OsPath []string
}

// generate the GOPATH environment pair
func (w *Workspace) EnvGOPATH() string {
	return WS_ENV + "=" + w.GenerateGOPATH()
}

// generate the GOBIN environment pair
func (w *Workspace) EnvGOBIN() string {
	return PKG_ENV + "=" + w.GenerateGOBIN()
}

// generate the GOPATH value
func (w *Workspace) GenerateGOPATH() string {
	return concatPath(w.GoPath, w.Root)
}

// generate the GOBIN value
func (w *Workspace) GenerateGOBIN() string {
	return path.Join(w.Root, BIN_DIR)
}

// generate the OS PATH value
func (w *Workspace) GeneratePATH() string {
	return concatPath(w.OsPath, w.GenerateGOBIN())
}

func (w *Workspace) String() string {
	return "Workspace(" + w.Root + ")"
}

func init() {
	WS_DEFAULT, _ = os.Getwd()
}

// create a workspace instance. the path array is used to determine
// the working directory, the place to install compiled binaries
// as well as additional lookup directories for core or external
// libraries. the values of _paths_ are expected to be absolute
// directory names.
// the the SDK directory is not an empty string, it is expected to
// contain the path to a GO installation which will be included in
// the OS path (its **bin** subdirectory to be precise).
// the _keepEnv_ directive will reuse the existing GOPATH directories
// and prepend the workspace directories.
func ParseWorkspace(paths []string, sdk string, keepEnv bool) (ws *Workspace, err error) {
	var gopath []string
	var ospath []string
	var langdir string
	var workdir string

	// generate GOPATH

	switch len(paths) {
	case 0:
		T("using", WS_DEFAULT, "as the workspace root")
		gopath = []string{}
		workdir = WS_DEFAULT
	case 1:
		T("workspace root:", paths[0], "no other includes")
		gopath = []string{}
		workdir = paths[0]
	default:
		T("workspace root:", paths[0], "+ includes")
		gopath = paths[1:]
		workdir = paths[0]
	}

	if keepEnv {
		D("appending", WS_ENV, "to workspace path")
		gopath = extendPath(gopath, WS_ENV)
	}

	// generate PATH

	if 0 == len(sdk) {
		ospath = extendPath([]string{}, OS_ENV)
	} else {
		D("using custom GO installation", sdk)
		langdir = path.Join(sdk, BIN_DIR)
		ospath = extendPath([]string{langdir}, OS_ENV)
	}

	return &Workspace{workdir, gopath, ospath}, nil
}

func concatPath(path []string, directory string) string {
	sep := string(os.PathListSeparator)
	postfix := strings.Join(path, sep)

	if 0 == len(path) {
		return directory
	} else {
		return directory + sep + postfix
	}
}

func extendPath(path []string, environment string) []string {
	extension := os.Getenv(environment)
	fragments := strings.Split(extension, string(os.PathListSeparator))

	return append(path, fragments...)
}
