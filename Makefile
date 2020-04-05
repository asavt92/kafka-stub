PROJECTNAME := $(shell basename "$(PWD)")
# Go related variables.
GOBASE := $(shell pwd)
GOPATH := $(GOBASE)/vendor:$(GOBASE)
GOBIN := $(GOBASE)

.PHONY: build
build:
	 CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o main ./cmd/main.go

.PHONY: docker
docker: build
	docker build -t $(PROJECTNAME) .

.PHONY: go-get
go-get:
	@echo "  >  Checking if there is any missing dependencies..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get $(get)

.DEFAULT_GOAL := build