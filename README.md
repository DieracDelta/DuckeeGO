# DuckDuckGo
Concolic Execution Engine for Golang

Relies upon aclements/z3 package.

## Setup stuff

### Downloading Z3
- get the latest release from `https://github.com/Z3Prover/z3/releases` (we are using 4.6.0)
- when `go get`-ing `github.com/aclements/go-z3/z3`, use `CGO_CFLAGS="-I/path/to/directory/with/z3.h"`
- also need to copy relevant libraries from the `include` directory into wherever your libraries are stored. (for macs, this is `/usr/local/lib`, for windows, this is `C:\\Windows\\System32`)

### Running our application

execute `bash src/run.sh`

Include path of config.json file as script paramter to change directory. See example directory for an example project to concolically execute. Upon execution of run.sh, DUCKEEGO will build instrumented code to allow for conolic execution inside of /src/tmp. The user may then cd into this directory, and run `go build && go run` to execute code concolically.
