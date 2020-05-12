#!/bin/bash 

BASEDIR=`dirname $0`/
BASEDIR=`(cd "$BASEDIR"; pwd)`

SRC_DIR="$BASEDIR/src"
SRC_FILE="$BASEDIR/src/totoroAgent.go"
DIST_DIR="$BASEDIR/bin"
DIST_FILE="$BASEDIR/bin/totoroAgent"

export GOPATH=$BASEDIR

if [ ! -d "$DIST_DIR" ];then
	mkdir -p $DIST_DIR
fi

echo "$BASEDIR"

cd $SRC_DIR

go clean
go install totoroAgent
go build -o $DIST_FILE $SRC_FILE

echo "build sucessful, dist file: $DIST_FILE"
