# gospace

shell spawner for go development workspaces

# synopsis

gospace \[OPTION\]... \[PATH\]...  

# description

in order to create a new workspace fast without a lot of manual scaffolding,
gospace spawns a shell with the appropriate environment variables set.

one might argue that this is against the golang conventions, but it allows
developing go projects outside your regular workspace (e.g. quick hacks)

it also eases the development process when using multiple installations
of Go (e.g. vanilla GO and _Google App Engine_ GO).

continuous integration can also use this script to compile an application
against multiple versions of GO.

----

gospace does:

1. define/overwrite/extend GOPATH
2. export GOBIN as the first entry of GOPATH + _bin_
2. extend PATH with GOBIN and optionally GOHOME/bin
3. spawn a shell
    * the binary path provided on the commandline
    * $SHELL
    * /bin/sh
4. set the working directory to the first entry of GOPATH

each path on the commandline is resolved against

* $PWD
* $GOSPACES
* $CDPATH

the behaviour of _gospace_ can be controlled via command-line arguments:

    -b, --blank           do not reuse GOPATH is defined
    -n, --dry             simulates the shell spawning
    -v, --verbose         raise the verbosity
    -s, --shell=PATH      use the provided shell in the workspace
    -g, --go=[DIR]        include DIR/bin or GOHOME/bin in the shell PATH
    -h, --help            display the usage message and exit
    -V, --version         print the gospace command version and exit

including -- on the commandline causes all remaining arguments to be passed
on to the shell command.

# environment

**CDPATH** and **GOSPACES** are both directory resolution inputs. they are
expected to contain colon separated directories. relative entries are resolved
against the current working directory.

**GOSPACES** exists soley for the purpose of having a set of project root
directories without interfering/clogging **CDPATH** (it is used by some
shells for switching directories). directories found via **GOSPACES** take
precedence over **CDPATH**.

if no shell has been defined **SHELL** is used. if this environment variable
does not exist as well, _/bin/sh_ is the gospace shell of choice.

specifying _--go_ on the commandline causes the evaluation of **GOHOME**,
unless the parameter has a value. **GOHOME** is expected to point to the
installation directory of a golang installation. its subdirectory _bin_
will be included in the PATH of the spawned shell.

# examples

> gospace

create/rewrite GOPATH to include the PWD as its first element

> gospace --blank

set GOPATH to PWD, regardless of the previous value

> gospace $HOME/goprojects/myproject

use the provided path as go workspace

> gospace myproject /usr/lib/go-contrib

search $CDPATH for a potential directory with that name and use it as
workspace. _/usr/lib/go-contrib_ will be included in the GOPATH.

> gospace --shell=/bin/bash myproject -- --login

use `/bin/bash` as the working shell instead of $SHELL. the command is invoked
with the additional _--login_ argument.

> gospace --go=/opt/go-gae/bin -- -c "go test"
