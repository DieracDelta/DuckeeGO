go build configStuff.go main.go
# see /tmp/DuckieConcolic for the stuff
rm -dr /tmp/DuckieConcolic
./main ../example/config.json
