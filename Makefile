MAKEFLAGS=-j1

.PHONY: build clean
build: clean; go build -o target/ffcorpus cmd/ffcorpus/main.go
clean: ; rm -rf target
