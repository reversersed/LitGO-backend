API_DIRECTORIES = api_gateway api_user

run: clean gen test start


install: i

i:
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install github.com/golang/mock/mockgen@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@$(MAKE) clean

gen:
	@echo Generating protobuf files...
	@cd ./api_gateway && protoc -I ../proto --go_out=. --go-grpc_out=. ../proto/*.proto && cd ..
	@cd ./api_user && protoc -I ../proto --go_out=. --go-grpc_out=. ../proto/*.proto && cd ..
	@echo Code generated successfully

	@swag init --parseDependency -d ./api_gateway/internal/handlers -g ../app/app.go -o ./api_gateway/docs
	@cd ./api_gateway/ && go generate ./...
	@cd ./api_user/ && go generate ./...

upgrade: clean
	@cd ./api_gateway/ && go get -u ./... && go mod tidy
	@cd ./api_user/ && go get -u ./... && go mod tidy

clean:
	@cd ./api_gateway/ && go mod tidy
	@cd ./api_user/ && go mod tidy

start:
	@docker compose up --build --timestamps --wait --wait-timeout 1800 --remove-orphans -d

stop:
	@docker compose stop

test-verbose:
	@cd ./api_gateway/ && go generate ./... && go test ./... -v
	@cd ./api_user/ && go generate ./... && go test ./... -v

test: test-folder-creation
	@cd ./api_gateway/ && go generate ./... && go test ./... -coverprofile=tests/coverage -coverpkg=./... && go tool cover -func=tests/coverage -o tests/coverage.func && go tool cover -html=tests/coverage -o tests/coverage.html
	@cd ./api_user/ && go generate ./... && go test ./... -coverprofile=tests/coverage -coverpkg=./... && go tool cover -func=tests/coverage -o tests/coverage.func && go tool cover -html=tests/coverage -o tests/coverage.html

test-folder-creation:
ifeq ($(OS),Windows_NT)
	-@cd ./api_gateway/ && mkdir tests
	-@cd ./api_user/ && mkdir tests
else
	-@cd ./api_gateway/ && mkdir -p tests
	-@cd ./api_user/ && mkdir -p tests
endif