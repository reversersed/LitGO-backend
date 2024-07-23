API_DIRECTORIES = api_gateway api_user

run: gen start

gen:
	@echo Generating protobuf files...
	@cd ./api_gateway && protoc -I ../proto --go_out=. --go-grpc_out=. ../proto/*.proto && cd ..
	@cd ./api_user && protoc -I ../proto --go_out=. --go-grpc_out=. ../proto/*.proto && cd ..
	@echo Code generated successfully

	@swag init --parseDependency -d ./api_gateway/internal/handlers -g ../app/app.go -o ./api_gateway/docs
	@cd ./api_gateway/ && go generate ./...
	@cd ./api_user/ && go generate ./...

start:
	@docker compose up --build --timestamps --wait --wait-timeout 1800 --remove-orphans -d

stop:
	@docker compose stop