package gospace

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var (
	// fallback shell path
	SHELL_DEFAULT string = "/bin/sh"
	// environment variable to read the shell path
	SHELL_ENV string = "SHELL"
	// environment variable to read for binary lookups
	PATH_ENV string = "PATH"
	// shell resolver error
	resolveError error = errors.New("Unable to find any suitable shell")
)

// command + arguments
type Shell struct {
	Path string
	Args []string
}

type builder struct {
	artifact string
}

// execute the shell. the shell will be invoked with the internal
// commandline arguments. the environment is enhanced with various
// go related variables. stdin, stdout and stderr are attached to
// the sub-process.
//
func (s *Shell) Launch(workspace *Workspace, simulate bool) error {
	var shell *exec.Cmd = exec.Command(s.Path, s.Args...)

	if err := os.Setenv(PATH_ENV, workspace.GeneratePATH()); nil != err {
		return err
	}

	shell.Dir = workspace.Root
	shell.Stdin = os.Stdin
	shell.Stdout = os.Stdout
	shell.Stderr = os.Stderr
	shell.Env = append(os.Environ(),
		workspace.EnvGOPATH(),
		workspace.EnvGOBIN())

	return shell.Run()
}

func (s *Shell) String() string {
	argv := ""

	if 0 < len(s.Args) {
		argv = " " + strings.Join(s.Args, " ")
	}

	return fmt.Sprintf("Shell(%s%s)", s.Path, argv)
}

func (b *builder) hasArtifact() bool {
	return 0 < len(b.artifact)
}

func (b *builder) use(path string) (self *builder) {
	self = b

	T("attempting to resolve shell", path)

	if b.hasArtifact() {
		D("shell already found; skipping lookup")
		return
	} else if 0 == len(path) {
		D("no shell path provided for lookup")
		return
	} else if abs, ok := SearchPathEnvironment(PATH_ENV, path); false == ok {
		D("shell lookup via PATH failed")
		return
	} else if DirExists(abs) {
		I("shell path is a directory")
		return
	} else {
		// TODO: check permission bits for executable flag
		b.artifact = abs

		return
	}
}

func (b *builder) build() (string, error) {
	if b.hasArtifact() {
		return b.artifact, nil
	}

	return "", resolveError
}

// resolve the path against various sources. the first match is used.
// if the path is empty, the shell specified in the environment as
// _SHELL_, otherwise the fallback value **/bin/sh** is used.
// if the path is not absolute it is resolved against the PATH
// directories. if all lookups yield no usable result, an error
// is returned.
func ResolveShell(path string, args []string) (*Shell, error) {
	builder := builder{""}
	osshell := os.Getenv(SHELL_ENV)

	builder.
		use(path).
		use(osshell).
		use(SHELL_DEFAULT)

	if binary, err := builder.build(); nil == err {
		D("using shell", binary)

		return &Shell{binary, args}, nil
	} else {
		return nil, err
	}
}
