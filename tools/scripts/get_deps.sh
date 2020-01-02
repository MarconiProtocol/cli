#!/bin/bash

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null && pwd)"

cd $DIR/../../
go get -d .
go get -u github.com/gorilla/rpc
