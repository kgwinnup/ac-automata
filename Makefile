
.PHONY: example test

test:
	go test -v ./...

example:
	cd example && go build -o ../aho main.go
