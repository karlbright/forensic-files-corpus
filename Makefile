MAKEFLAGS=-j1

.PHONY: setup build clean lambda cleanl release
setup: ; mkdir -p target release
copy: ; cp sentences.txt target/sentences.txt
build: setup cleanc; go build -o target/ffcorpus cmd/ffcorpus/*.go && make copy
clean: ; rm -f target/ffcorpus
lambda: setup cleanl; GOOS=linux go build -o target/skippalenik cmd/skippalenik/main.go && make copy
cleanl: ; rm -f target/skippalenik && rm -f release/lambda.zip
release: lambda; zip release/skippalenik.zip target/skippalenik target/sentences.txt
