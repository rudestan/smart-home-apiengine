build_all:
	go build -race -o ./build/smh-webserver ./cmd/webserver/serve.go
	go build -race -o ./build/smh-configurator ./cmd/configurator/configure.go
	go build -race -o ./build/smh-runner ./cmd/runner/run.go

test:
	go test -race ./...

coverage:
	go test -race -v -coverprofile ./build/cover.out ./...
	go tool cover -html=./build/cover.out -o ./build/cover.html