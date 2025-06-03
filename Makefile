build:
	go build -o paw

test:
	go test -v ./tests/...

release:
	goreleaser release --clean

release_check:
	goreleaser check
