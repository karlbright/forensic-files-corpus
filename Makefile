MAKEFLAGS=-j1

.PHONY: build clean
clean: ; rm -rf target
build: clean; go build -o target/ffcorpus cmd/ffcorpus/main.go
