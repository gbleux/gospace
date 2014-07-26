package flag

// commandline arguments registry.
// contains the parsed values from the commandline
type Arguments struct {
	Blank     bool
	NoRun     bool
	GoSDK     string
	Shell     string
	ShellArgv []string
	Path      []string
}

// convenient wrapper to append a value to the shell argument slice
func (a *Arguments) AppendShellArgument(value string) {
	a.ShellArgv = append(a.ShellArgv, value)
}

// convenient wrapper to append a value to the path slice
func (a *Arguments) AppendPath(value string) {
	a.Path = append(a.Path, value)
}

// argument registry instance factory
func NewArguments() *Arguments {
	shellParams := []string{}
	includePath := []string{}

	return &Arguments{false, false, "", "", shellParams, includePath}
}
