include ../Makefile.base

IMAGE_NAME=harvey/switchboard
VERSION=1.0

BUILD_FILE=root/usr/bin/switchboard
SRC_ENTRYPOINT=src/switchboard.go
PLUGINS := $(shell ls -d plugin/*/*)
PLUGIN_DIR=root/usr/lib

plugins:
	for plugin in $(PLUGINS); do \
		GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o $(PLUGIN_DIR)/switchboard_$$(echo $$plugin | tr / _).so $$plugin/main.go; \
	done

build:
	go get -d ./...
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_FILE) $(SRC_ENTRYPOINT)
	docker build -t $(IMAGE_NAME_TAGGED) -t $(IMAGE_NAME_VERSION) .

run:
	docker run \
		--name $(NAME) \
		-p 10000:80 \
		--rm \
		-v ~/.aws:/root/.aws \
		-v $(CURDIR)/etc/raw.yml:/etc/switchboard/conf.yml:ro \
		-v /tmp/switchboard:/tmp/switchboard \
		$(IMAGE_NAME_VERSION)
