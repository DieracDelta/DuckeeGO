DEST=$1
if [ $# -eq 0 ]
  then
  DEST="../example/config.json"
fi
go build main.go jsonDefs.go instrumentationHelpers.go addInstrumentation.go && rm -dr tmp/DuckieConcolic ; ./main $DEST
