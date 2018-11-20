.PHONY: build-docker

DOCKER_TAG ?= latest

build-docker: $(shell find ./ -name "*.go") Dockerfile
	docker build -t klardotsh/shuffleboard:$(DOCKER_TAG) .
