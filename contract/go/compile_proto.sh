#!/bin/bash

set -eux

protoc --proto_path=$GOPATH/src/github.com/BeDreamCoder/uwavm --go_out=$GOPATH/src contract/pb/contract.proto
