build:
	go build

test:
	go test -v -count=1 ./...

cover:
	go test -v -count=1 -cover ./...

lint:
	golangci-lint run --enable-all

clean:
	go clean -modcache && go mod tidy

run:
	go run build/darwin/amd64/sandbox-api

compile:
	go-executable-build.sh sandbox-api .

