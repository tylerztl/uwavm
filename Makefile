# Copyright zhangtailin All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
# -------------------------------------------------------------
# This makefile defines the following targets
#
#   - protos - generate all protobuf artifacts based on .proto files
#   - networkUp - start the fabric network
#   - networkDown - teardown the fabric network and clean the containers and intermediate images
#   - satrt - start the fabric-sdk-go server

.PHONY: protos
protos :
	./contract/pb/compile.sh

.PHONY: build
build :
	go build -o output/uwavm run/main.go

.PHONY: wasm
wasm :
	GOOS=js GOARCH=wasm go build -o output/wasm/erc20.wasm example/erc20/erc20.go

.PHONY: clean
clean :
	rm -rf output
