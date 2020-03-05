build_all:
	go build -o ./bin/smh-webserver ./cmd/webserver/serve.go
	go build -o ./bin/smh-configurator ./cmd/configurator/configure.go
	go build -o ./bin/smh-runner ./cmd/runner/run.go

test:
	go test ./...

coverage:
	go test -v -coverprofile ./bin/cover.out ./...
	go tool cover -html=./bin/cover.out -o ./bin/cover.html