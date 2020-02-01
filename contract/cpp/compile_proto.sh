#!/bin/bash

set -e

protoc --proto_path=$GOPATH/src/github.com/BeDreamCoder/uwavm/contract/pb --cpp_out=pb contract.proto
