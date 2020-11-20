BUILD_DIR=build
CMD_DIR=cmd
CMDS=$(patsubst $(CMD_DIR)/%,%,$(wildcard $(CMD_DIR)/*))

all: fmt check 

fmt:
	gofmt -s -w -l .
	@goimports -w -local gobe $(shell find . -type f -name '*.go' -not -path "./internal/*")

check:
	golint -set_exit_status ./... && \
	go vet -all ./... && \
	misspell -error */** && \
	go mod tidy
