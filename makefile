build_all:
	go build -o ./build/smh-webserver ./cmd/webserver/serve.go
	go build -o ./build/smh-configurator ./cmd/configurator/configure.go
	go build -o ./build/smh-runner ./cmd/runner/run.go

test:
	go test ./...

coverage:
	go test -v -coverprofile ./build/cover.out ./...
	go tool cover -html=./build/cover.out -o ./build/cover.html