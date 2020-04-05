PROJECTNAME := $(shell basename "$(PWD)")


.PHONY: build
build:
	 CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o main ./cmd/main.go

.PHONY: docker
docker: build
	docker build -t $(PROJECTNAME) .

#.PHONY: test
#test:
#	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build