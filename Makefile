# Config
dockerHubOrg := "threethink"
imageName := "slack-notifier"
commitHash := $(shell git rev-parse HEAD)
fullImageName := $(dockerHubOrg)/$(imageName):$(commitHash)

all: build

build:
	@docker build --pull --rm -t $(fullImageName) .

push:
	@docker push $(fullImageName)

test:
	@docker run -p 8080:8080 -e NOTIFIER_WEBHOOKURL=$(webHookUrl) $(fullImageName)

minio:
	@docker run -d --rm --privileged -p 9000:9000 --name minio1 -v $(PWD):/data -v $(PWD):/root/.minio minio/minio server /data