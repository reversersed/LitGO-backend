API_DIRECTORIES = api_gateway api_user api_author api_genre api_book
PROTO_PKG_FOLDERS = users genres authors books
CMDSEP = &&

run: i clean gen test-verbose check start

install: i

i:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/golang/mock/mockgen@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@$(MAKE) clean

check:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@$(foreach directory,$(API_DIRECTORIES),\
	cd ./$(directory)/ && echo checking $(directory)... && golangci-lint run && cd .. \
	$(CMDSEP)) echo lint checks completed

fix:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@$(foreach directory,$(API_DIRECTORIES),\
	cd ./$(directory)/ && echo checking $(directory)... && golangci-lint run --fix && cd .. \
	$(CMDSEP)) echo lint check and fix completed

gen:
	@$(foreach directory,$(API_DIRECTORIES),\
		@cd ./$(directory) && go generate ./... && cd ..\
		$(CMDSEP)) echo go files generated

	@swag init --parseDependency -d ./api_gateway/internal/handlers -g ../app/app.go -o ./api_gateway/docs

upgrade: clean i
	@$(foreach directory,$(API_DIRECTORIES),\
		@cd ./$(directory) && go get -v -u ./... && go mod tidy && cd ..\
		$(CMDSEP)) echo packages upgraded

clean:
	@go env -w GOPROXY=direct
	@$(foreach directory,$(API_DIRECTORIES),\
		cd ./$(directory) && go get -v -u github.com/reversersed/LitGO-proto/gen/go@latest && go mod tidy && cd ..\
		$(CMDSEP)) echo mod files cleaned
	@go env -w GOPROXY=https://proxy.golang.org,direct

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
		cd ./$(directory) && go test ./... -coverprofile=../../data/tests/$(directory)/coverage -coverpkg=./... && go tool cover -func=../../data/tests/$(directory)/coverage -o ../../data/tests/$(directory)/coverage.func && go tool cover -html=../../data/tests/$(directory)/coverage -o ../../data/tests/$(directory)/coverage.html && cd ..\
		$(CMDSEP)) echo tests completed successfully

test-folder-creation:
ifeq ($(OS),Windows_NT)
	-@cd .. & mkdir data
	-@cd ../data & mkdir tests
	-@cd ../data/tests & $(foreach directory,$(API_DIRECTORIES),\
		mkdir $(directory)\
		$(CMDSEP)) echo test directories has been created
else
	-@$(foreach directory,$(API_DIRECTORIES),\
		cd ./$(directory) & mkdir tests -p & cd ..\
		$(CMDSEP)) echo test directories has been created
endif