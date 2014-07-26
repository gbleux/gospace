package main

import (
	"fmt"
	"os"
	"path"

	"gospace"

	"cli/gospace/flag"
)

const (
	HEADLINE = "shell spawner for go development workspaces"
	FOOTER   = `commandline parsing can be terminated using --. all remaining values
will be passed to the shell command.
if no PATH is specified, the current working directory is used.`
)

var commandline *flag.Parser
var binaryname string

func init() {
	resolver := flag.PathResolver(resolverProxy)

	commandline = flag.NewParser(gospace.SDK_ENV, gospace.SHELL_ENV, &resolver)
	binaryname = path.Base(os.Args[0])
}

func main() {
	var err error = nil
	var code int = 0

	help := flag.Callback(printHelp)
	version := flag.Callback(printVersion)
	workspace := flag.Callback(launchWorkspace)

	commandline.
		On(flag.ACTION_HELP, &help).
		On(flag.ACTION_VERSION, &version).
		On(flag.ACTION_GOSPACE, &workspace)

	if code, err = commandline.Parse(os.Args[1:]); nil != err {
		fmt.Println(err.Error())
	}

	os.Exit(code)
}

func resolverProxy(path string) (string, error) {
	return gospace.ResolveGospace(path)
}

func printHelp(params *flag.Arguments) (int, error) {
	flag.WriteUsage(os.Stdout, binaryname, HEADLINE, FOOTER)

	return 0, nil
}

func printVersion(params *flag.Arguments) (int, error) {
	fmt.Println(binaryname, gospace.VERSION)

	return 0, nil
}

func launchWorkspace(params *flag.Arguments) (int, error) {
	var sh *gospace.Shell
	var ws *gospace.Workspace
	var err error

	if sh, err = gospace.ResolveShell(params.Shell, params.ShellArgv); nil != err {
		return 1, err
	} else if ws, err = gospace.ParseWorkspace(params.Path, params.GoSDK, !params.Blank); nil != err {
		return 2, err
	} else if err = sh.Launch(ws, params.NoRun); nil != err {
		return 4, err
	}

	return 0, nil
}
