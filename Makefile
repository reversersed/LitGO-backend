API_DIRECTORIES = api_gateway api_user api_author api_genre api_book api_review api_collection
PROTO_PKG_FOLDERS = users genres authors books reviews collections
CMDSEP = &&

run: i clean gen test-verbose check start

install: i

i:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/golang/mock/mockgen@latest
	@$(MAKE) clean

check:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@golangci-lint cache clean
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
ifeq ($(OS),Windows_NT)
	@$(foreach directory,$(API_DIRECTORIES),\
		cd ./$(directory) && go test -v -short ./... | findstr /V mocks && cd ..\
		$(CMDSEP)) echo tests completed successfully
else
	@$(foreach directory,$(API_DIRECTORIES),\
		cd ./$(directory) && go test -v -short ./... | grep -v mocks && cd ..\
		$(CMDSEP)) echo tests completed successfully
endif

test-verbose:
ifeq ($(OS),Windows_NT)
	@$(foreach directory,$(API_DIRECTORIES),\
		cd ./$(directory) && go test -v ./... | findstr /V mocks && cd ..\
		$(CMDSEP)) echo tests completed successfully
else
	@$(foreach directory,$(API_DIRECTORIES),\
		cd ./$(directory) && go test -v ./... | grep -v mocks && cd ..\
		$(CMDSEP)) echo tests completed successfully
endif

test: test-folder-creation gen
ifeq ($(OS),Windows_NT)
	@$(foreach directory,$(API_DIRECTORIES),\
		cd ./$(directory) && go test -coverprofile=../../data/tests/$(directory)/coverage -coverpkg=./... ./... -json | go-test-report -o ../../data/tests/$(directory)/results.html | findstr /V mocks && go tool cover -func=../../data/tests/$(directory)/coverage -o ../../data/tests/$(directory)/coverage.func && go tool cover -html=../../data/tests/$(directory)/coverage -o ../../data/tests/$(directory)/coverage.html && cd ..\
		$(CMDSEP)) echo tests completed successfully
else
	@$(foreach directory,$(API_DIRECTORIES),\
		cd ./$(directory) && go test -coverprofile=../../data/tests/$(directory)/coverage -coverpkg=./... ./... | grep -v mocks && go tool cover -func=../../data/tests/$(directory)/coverage -o ../../data/tests/$(directory)/coverage.func && go tool cover -html=../../data/tests/$(directory)/coverage -o ../../data/tests/$(directory)/coverage.html && cd ..\
		$(CMDSEP)) echo tests completed successfully
endif

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