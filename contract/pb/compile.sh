#!/bin/bash

set -eux

protoc --proto_path=. --go_out=$GOPATH/src ./contract.proto
