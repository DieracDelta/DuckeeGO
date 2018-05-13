go build main.go jsonDefs.go instrumentationHelpers.go addInstrumentation.go && rm -dr tmp/DuckieConcolic && ./main ../example/config.json
