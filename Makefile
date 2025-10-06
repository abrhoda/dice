# optional flags to put into build task
PROD_BUILD_FLAGS ?= -ldflags="-s -w" -trimpath
DEBUG_BUILD_FLAGS ?= -gcflags="all=-N -l"

PACKAGE ?= ./...
EXECUTABLE_NAME ?= dice
OUT_DIR ?= ./out
COVER_PROFILE ?= coverage.out

.PHONY: build
build: ## buid the project
	@mkdir -p $(OUT_DIR)
	@go build -o $(OUT_DIR)/$(EXECUTABLE_NAME) $(PACKAGE)

.PHONY: test
test: ## run all test 
	@go test $(PACKAGE) -v

.PHONY: format
format: ## format project
	@go fmt ./...

.PHONY: cover
cover: ## run tests and generate the converage.out file
	@mkdir -p $(OUT_DIR)
	@go test $(PACKAGE) -coverprofile=$(OUT_DIR)/$(COVER_PROFILE)

.PHONY: coverhtml
coverhtml: cover ## generate the html coverage report to view at ./coverage.html
	@go tool cover -html=$(OUT_DIR)/$(COVER_PROFILE) -o $(OUT_DIR)/coverage.html

.PHONY: coverfunc
coverfunc: cover ## generate a report about % coverage by function to stdout
	@go tool cover -func=$(OUT_DIR)/$(COVER_PROFILE)

.PHONY: coverpercent
coverpercant: cover ## outputs the total unit test coverage % for project.
	@TOTAL=$$(@go tool cover -func=$(OUT_DIR)/$(COVER_PROFILE) | grep total)
	@echo TOTAL

.PHONY: cleancoverage
cleancoverage: ## clean up the generated coverage.out file
	@rm -f $(OUT_DIR)/$(COVER_PROFILE) $(OUT_DIR)/coverage.html

.PHONY: clean
clean: ## remove built files and dirs
	@rm -rf $(OUT_DIR)

.PHONY: vet
vet: ## runs a `go vet` check for the project
	@go vet
help: ## print this help message
	@awk -F ':|##' '/^[^\t].+?:.*?##/ {printf "\033[36m%-20s\033[0m %s\n", $$1, $$NF}' $(MAKEFILE_LIST)
