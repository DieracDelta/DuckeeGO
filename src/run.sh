go build main.go jsonDefs.go instrumentationHelpers.go addInstrumentation.go
# see /tmp/DuckieConcolic for the stuff
rm -dr /tmp/DuckieConcolic
./main ../example/config.json
