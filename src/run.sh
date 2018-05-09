go build main.go configStuff.go
# see /tmp/DuckieConcolic for the stuff
rm -dr /tmp/DuckieConcolic
./main ../example/config.json
