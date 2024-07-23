API_DIRECTORIES = api_gateway api_user

run: gen start

gen:
	@echo Generating protobuf files...
	@cd ./api_gateway && protoc -I ../proto --go_out=. --go-grpc_out=. ../proto/*.proto && cd ..
	@cd ./api_user && protoc -I ../proto --go_out=. --go-grpc_out=. ../proto/*.proto && cd ..
	@echo Code generated successfully

start:
	@docker compose up --build --timestamps --wait --wait-timeout 1800 --remove-orphans -d