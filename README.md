# DuckDuckGo
Concolic Execution Engine for Golang

relies upon https://github.com/fatih/astrewrite for augmenting ast while inserting instrumentation.


QUACK. quack quack quack.


## Setup stuff

### Downloading Z3
- get the latest release from `https://github.com/Z3Prover/z3/releases` (we are using 4.6.0)
- when `go get`-ing `github.com/aclements/go-z3/z3`, use `CGO_CFLAGS="-I/path/to/directory/with/z3.h"`
- also need to copy relevant libraries from the `include` directory into wherever your libraries are stored. (for macs, this is `/usr/local/lib`, for windows, this is `C:\\Windows\\System32`)

### Other stuff
- run `go get https://github.com/fatih/astrewrite`

