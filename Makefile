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
	cd ./$(directory)/ $(CMDSEP) echo checking $(directory)... $(CMDSEP) golangci-lint run $(CMDSEP) cd .. \
	$(CMDSEP)) cd ./pkg $(CMDSEP) echo checking $(directory)... $(CMDSEP) golangci-lint run

fix:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@$(foreach directory,$(API_DIRECTORIES),\
	cd ./$(directory) $(CMDSEP) echo checking $(directory)... $(CMDSEP) golangci-lint run --fix $(CMDSEP) cd .. \
	$(CMDSEP)) cd ./pkg $(CMDSEP) echo checking $(directory)... $(CMDSEP) golangci-lint run --fix

gen:
	@$(foreach directory,$(API_DIRECTORIES),\
		@cd ./$(directory) $(CMDSEP) go generate ./... $(CMDSEP) cd ..\
		$(CMDSEP)) cd ./pkg $(CMDSEP) go generate ./...

upgrade: clean i
	@$(foreach directory,$(API_DIRECTORIES),\
		@cd ./$(directory) $(CMDSEP) go get -v -u ./... $(CMDSEP) go mod tidy $(CMDSEP) cd ..\
		$(CMDSEP)) cd ./pkg $(CMDSEP) go get -v -u ./...

clean:
	@go env -w GOPROXY=direct
	@$(foreach directory,$(API_DIRECTORIES),\
		cd ./$(directory) $(CMDSEP) go get -v -u github.com/reversersed/LitGO-proto/gen/go@latest $(CMDSEP)\
		go get -v -u github.com/reversersed/LitGO-backend-pkg@latest $(CMDSEP)\
		go mod tidy $(CMDSEP) cd ..\
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
		cd ./$(directory) $(CMDSEP) go test -v -short ./... | findstr /V mocks $(CMDSEP) cd ..\
		$(CMDSEP))cd ./pkg $(CMDSEP) go test -v -short ./... | findstr /V mocks
else
	@$(foreach directory,$(API_DIRECTORIES),\
		cd ./$(directory) $(CMDSEP) go test -v -short ./... | grep -v mocks $(CMDSEP) cd ..\
		$(CMDSEP))cd ./pkg $(CMDSEP) go test -v -short ./... | grep -v mocks
endif

test-verbose:
ifeq ($(OS),Windows_NT)
	@$(foreach directory,$(API_DIRECTORIES),\
		cd ./$(directory) $(CMDSEP) go test -v ./... | findstr /V mocks $(CMDSEP) cd ..\
		$(CMDSEP)) cd ./pkg $(CMDSEP) go test -v ./... | findstr /V mocks
else
	@$(foreach directory,$(API_DIRECTORIES),\
		cd ./$(directory) $(CMDSEP) go test -v ./... | grep -v mocks $(CMDSEP) cd ..\
		$(CMDSEP)) cd ./pkg $(CMDSEP) go test -v ./... | grep -v mocks
endif

test: test-folder-creation gen
ifeq ($(OS),Windows_NT)
	@$(foreach directory,$(API_DIRECTORIES),\
		cd ./$(directory) $(CMDSEP) go test -coverprofile=../../data/tests/$(directory)/coverage -coverpkg=./... ./... -json | go-test-report -o ../../data/tests/$(directory)/results.html | findstr /V mocks $(CMDSEP) go tool cover -func=../../data/tests/$(directory)/coverage -o ../../data/tests/$(directory)/coverage.func $(CMDSEP) go tool cover -html=../../data/tests/$(directory)/coverage -o ../../data/tests/$(directory)/coverage.html $(CMDSEP) cd ..\
		$(CMDSEP)) cd ./pkg $(CMDSEP) go test -coverprofile=../../data/tests/pkg/coverage -coverpkg=./... ./... -json | go-test-report -o ../../data/tests/pkg/results.html | findstr /V mocks $(CMDSEP) go tool cover -func=../../data/tests/pkg/coverage -o ../../data/tests/pkg/coverage.func $(CMDSEP) go tool cover -html=../../data/tests/pkg/coverage -o ../../data/tests/pkg/coverage.html
else
	@$(foreach directory,$(API_DIRECTORIES),\
		cd ./$(directory) $(CMDSEP) go test -coverprofile=../../data/tests/$(directory)/coverage -coverpkg=./... ./... | grep -v mocks $(CMDSEP) go tool cover -func=../../data/tests/$(directory)/coverage -o ../../data/tests/$(directory)/coverage.func $(CMDSEP) go tool cover -html=../../data/tests/$(directory)/coverage -o ../../data/tests/$(directory)/coverage.html $(CMDSEP) cd ..\
		$(CMDSEP)) cd ./pkg $(CMDSEP) go test -coverprofile=../../data/tests/pkg/coverage -coverpkg=./... ./... | grep -v mocks $(CMDSEP) go tool cover -func=../../data/tests/pkg/coverage -o ../../data/tests/pkg/coverage.func $(CMDSEP) go tool cover -html=../../data/tests/pkg/coverage -o ../../data/tests/pkg/coverage.html
endif

test-folder-creation:
ifeq ($(OS),Windows_NT)
	-@cd .. & mkdir data
	-@cd ../data & mkdir tests
	-@cd ../data/tests & $(foreach directory,$(API_DIRECTORIES),\
		mkdir $(directory)\
		$(CMDSEP)) mkdir pkg
else
	-@cd .. & mkdir data -p
	-@cd ../data & mkdir tests -p
	-@cd ../data/tests & $(foreach directory,$(API_DIRECTORIES),\
		mkdir $(directory) -p\
		$(CMDSEP)) mkdir pkg -p
endif