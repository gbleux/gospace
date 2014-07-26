package flag

import (
	"fmt"
	"io"
	"strings"

	"gospace"
)

const (
	// trigger the _help_ action
	ACTION_HELP Action = iota
	// trigger the _version_ action
	ACTION_VERSION = iota
	// trigger the _gospace_ action
	ACTION_GOSPACE = iota
)

var (
	version *Parameter
	help    *Parameter
	norun   *Parameter
	blank   *Parameter
	shell   *Parameter
	gosdk   *Parameter
	debug   *Parameter
)

// typedef for triggers
type Action int

// action callback
type Callback func(argv *Arguments) (int, error)

// conversion utility to resolve a (possibly) relative path
type PathResolver func(path string) (string, error)

// command-line parser
type Parser struct {
	gosdkEnv  string
	shellEnv  string
	resolver  *PathResolver
	callbacks map[Action]*Callback
}

// register a callback to be invoked when the commandline contains
// certain action triggers. if the commandline contains multiple
// triggers, the first encountered (not registered) will be executed.
func (p *Parser) On(action Action, cb *Callback) (self *Parser) {
	self = p

	gospace.T("registering callback", action)

	p.callbacks[action] = cb

	return
}

// process the arguments and call the first matching trigger callback.
// the return values are most likely from the callback itself, unless
// an unresolvable directory was provided on the commandline.
func (p *Parser) Parse(input []string) (status int, err error) {
	var passthrough bool = false
	var argv *Arguments = NewArguments()

	gospace.T("processing commandline", input)

	for _, arg := range input {
		if passthrough {
			gospace.T("found argument for sub-shell")

			argv.AppendShellArgument(arg)
		} else {
			gospace.T("processing gospace argument", arg)

			switch {
			case "--" == arg:
				gospace.T("argument terminator encountered")
				passthrough = true
			case help.Matches(arg):
				gospace.T("help action triggered")
				return p.fire(ACTION_HELP, argv)
			case version.Matches(arg):
				gospace.T("version action triggered")
				return p.fire(ACTION_VERSION, argv)
			case norun.Matches(arg):
				gospace.T("shell spawning is only simulated")
				argv.NoRun = true
			case blank.Matches(arg):
				gospace.T("blank flag defined")
				argv.Blank = true
			case debug.Matches(arg):
				gospace.T("increasing verbosity")
				// increase by 2, so -vvv will yield full verbosity
				gospace.LOG_LEVEL.Increase(2)
			case shell.Matches(arg):
				gospace.T("custom shell argument")
				argv.Shell = shell.ParseValueOr(arg, p.shellEnv, "")
			case gosdk.Matches(arg):
				gospace.T("custom go installation provided")
				argv.GoSDK = gosdk.ParseValueOr(arg, p.gosdkEnv, "")
			case strings.HasPrefix(arg, "-"):
				return 0, fmt.Errorf("Unknown argument '%s'", arg)
			default:
				gospace.T("received directory input for GOPATH")
				if path, err := (*p.resolver)(arg); nil != err {
					return 0, err
				} else {
					argv.AppendPath(path)
				}
			}
		}
	}

	return p.fire(ACTION_GOSPACE, argv)
}

func (p *Parser) fire(action Action, argv *Arguments) (int, error) {
	if callback, ok := p.callbacks[action]; ok {
		return (*callback)(argv)
	}

	gospace.T("no callback registered for", action)

	return 0, nil
}

func init() {
	version = NewFlagParameter('V', "version", "display the application version and exit")
	help = NewFlagParameter('h', "help", "show this message and exit")
	norun = NewFlagParameter('n', "dry", "simulates the shell spawning")
	blank = NewFlagParameter('b', "blank", "overwrite GOPATH instead of extending it")
	debug = NewFlagParameter('v', "verbose", "raise the verbosity")
	shell = NewArgParameter('s', "shell", "PATH", "run the workspace in a custom shell")
	gosdk = NewArgParameter('g', "go", "PATH", "include the go installation in the PATH")
}

// parser instance factory
func NewParser(sdk string, shell string, resolver *PathResolver) *Parser {
	callbacks := make(map[Action]*Callback)

	return &Parser{sdk, shell, resolver, callbacks}
}

// write the program header, footer, usage and commandline
// arguments to the writer.
func WriteUsage(out io.Writer, application string, description string, footer string) {
	header := fmt.Sprintf("usage: %s [OPTION]... [PATH]...\n", application)

	io.WriteString(out, header)
	io.WriteString(out, description)
	io.WriteString(out, "\n\n")

	io.WriteString(out, "arguments:\n")
	io.WriteString(out, blank.Usage())
	io.WriteString(out, norun.Usage())
	io.WriteString(out, debug.Usage())
	io.WriteString(out, gosdk.Usage())
	io.WriteString(out, shell.Usage())
	io.WriteString(out, help.Usage())
	io.WriteString(out, version.Usage())
	io.WriteString(out, footer)
	io.WriteString(out, "\n")
}
