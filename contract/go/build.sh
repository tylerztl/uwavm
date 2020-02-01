#!/bin/bash

GOOS=js GOARCH=wasm go build -o erc20.wasm example/erc20.go
