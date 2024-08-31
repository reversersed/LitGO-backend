API_DIRECTORIES = api_gateway api_user api_author api_genre
PROTO_PKG_FOLDERS = users genres authors
CMDSEP = &&

run: clean gen test-verbose check start

install: i

i:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/favadi/protoc-go-inject-tag@latest
	@go install github.com/golang/mock/mockgen@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@$(MAKE) clean

check:
	@$(foreach directory,$(API_DIRECTORIES),\
	cd ./$(directory)/ && echo checking $(directory)... && golangci-lint run && cd .. \
	$(CMDSEP)) echo lint checks completed

fix:
	@$(foreach directory,$(API_DIRECTORIES),\
	cd ./$(directory)/ && echo checking $(directory)... && golangci-lint run --fix && cd .. \
	$(CMDSEP)) echo lint check and fix completed

gen:
	@$(foreach directory,$(API_DIRECTORIES),\
		cd ./$(directory)/ && protoc -I ../proto --go_out=. --go-grpc_out=. ../proto/*.proto && $(foreach folder,$(PROTO_PKG_FOLDERS),\
		protoc-go-inject-tag -input="./pkg/proto/$(folder)/*.pb.go" -remove_tag_comment &&)\
		protoc-go-inject-tag -input="./pkg/proto/*.pb.go" -remove_tag_comment && cd .. \
		$(CMDSEP)) echo proto files generated

	@swag init --parseDependency -d ./api_gateway/internal/handlers -g ../app/app.go -o ./api_gateway/docs

	@$(foreach directory,$(API_DIRECTORIES),\
		@cd ./$(directory) && go generate ./... && cd ..\
		$(CMDSEP)) echo go files generated

upgrade: clean i
	@$(foreach directory,$(API_DIRECTORIES),\
		@cd ./$(directory) && go get -u ./... && go mod tidy && cd ..\
		$(CMDSEP)) echo packages upgraded

clean:
	@$(foreach directory,$(API_DIRECTORIES),\
		cd ./$(directory) && go mod tidy && cd ..\
		$(CMDSEP)) echo mod files cleaned

start:
	@docker compose up --build --timestamps --wait --wait-timeout 1800 --remove-orphans -d

stop:
	@docker compose stop

up:
	@docker compose up --timestamps --wait --wait-timeout 1800 --remove-orphans -d

test-unit:
	@$(foreach directory,$(API_DIRECTORIES),\
		cd ./$(directory) && go test ./... -v -short && cd ..\
		$(CMDSEP)) echo tests completed successfully

test-verbose:
	@$(foreach directory,$(API_DIRECTORIES),\
		cd ./$(directory) && go test ./... -v && cd ..\
		$(CMDSEP)) echo tests completed successfully

test: test-folder-creation gen
	@$(foreach directory,$(API_DIRECTORIES),\
		cd ./$(directory) && go test ./... -coverprofile=tests/coverage -coverpkg=./... && go tool cover -func=tests/coverage -o tests/coverage.func && go tool cover -html=tests/coverage -o tests/coverage.html && cd ..\
		$(CMDSEP)) echo tests completed successfully

test-folder-creation:
ifeq ($(OS),Windows_NT)
	-@$(foreach directory,$(API_DIRECTORIES),\
		cd ./$(directory) & mkdir tests & cd ..\
		$(CMDSEP)) echo test directories has been created
else
	-@$(foreach directory,$(API_DIRECTORIES),\
		cd ./$(directory) & mkdir tests -p & cd ..\
		$(CMDSEP)) echo test directories has been created
endif