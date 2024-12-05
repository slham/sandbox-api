build:
	go build

build_cover:
	go build -cover

test: unit int

cover: unit_cover int_cover

unit:
	go test -tags=unit -count=1 ./...

unit_cover:
	go test -v -count=1 -tags=unit -cover ./...

int:
	go test -tags=integration -count=1 ./...

int_cover:
	go test -v -count=1 -tags=integration -cover ./...

int_cover_see:
	go tool covdata percent -i=./integration

lint:
	golangci-lint run --enable-all

lint_fix:
	golangci-lint run --enable-all --fix

clean:
	go clean -modcache && go mod tidy

run:
	./build/darwin/amd64/sandbox-api

run_cover:
	make build_cover && GOCOVERDIR=integration ./sandbox-api

compile:
	go-executable-build.sh sandbox-api .

build_image:
	docker build --rm -t "sandbox-api" --build-arg ex_path=linux/amd64 .

containerize: compile build_image

run_container:
	docker run -d --env-file ./env/local.env sandbox-api
