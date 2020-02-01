
.PHONY: build
build :
	go build -o output/uwavm run/main.go

.PHONY: clean
clean :
	rm -rf output
