build:
	@go build -C ./src/cmd/duwa -o ../../../bin/duwa

build-wasm:
	set TEMP=U:\projects\skybox\duwa\chewa\temp
	@GOOS=js GOARCH=wasm tinygo build -o ../duwa-site/public/duwa.wasm -opt 1 ./src/cmd/wasm/duwa.go

build-docker:
	@docker-compose up -d
	@docker-compose exec tinygo-dev tinygo build -o ./bin/duwa.wasm -opt 1 -target=wasm ./src/cmd/wasm/duwa.go
	@docker-compose exec tinygo-dev cp /tinygo/targets/wasm_exec.js /app/bin/
	@cp ./bin/duwa.wasm ../duwa-site/public/duwa.wasm
	@cp ./bin/wasm_exec.js ../duwa-site/public/wasm_exec.js
	@docker-compose stop

build-all: build build-wasm

build-tools:
	@go build -C src/cmd/tools -o ../../../bin/duwa-tools.exe

dev:
	@go run ./src/cmd/duwa/duwa.go

run: build
	@./bin/duwa

lint:
	@golangci-lint run -c ./golangci.yml ./...

test:
	@go test ./... -v --cover

test-report:
	@go test ./... -v --cover -coverprofile=coverage.out
	@go tool cover -html=coverage.out


docs-json:
	@go run src/cmd/tools/* gen docs

gen-std:
	@go run src/cmd/tools/* gen std-types


lsp:
	@go run src/cmd/tools/* lsp -v -p websocket