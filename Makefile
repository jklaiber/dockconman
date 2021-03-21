REGISTRY_NAME?=docker.io/jklaiber
IMAGE_VERSION?=0.0.0

.PHONY: all dev-start dev-up dev-down build container push clean

all: build

dev-start: dev-up
		go run ./cmd/dockconman -n dockconman-dev -c bash -d true

dev-up:
		docker run -d -it --name dockconman-dev ubuntu:focal bash

dev-down:
		docker rm -f dockconman-dev

build:
		mkdir -p bin
		$(MAKE) -C ./cmd/dockconman compile-dockconman

container:
		docker build -t $(REGISTRY_NAME)/dockconman:$(IMAGE_VERSION) -f ./Dockerfile .

push: container
		docker push $(REGISTRY_NAME)/dockconman:$(IMAGE_VERSION)

clean:
		rm -rf bin