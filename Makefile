# Config
dockerHubOrg := "threethink"
imageName := "slack-notifier"
commitHash := $(shell git rev-parse HEAD)
fullImageName := $(dockerHubOrg)/$(imageName):$(commitHash)

all: build

build:
	@docker build --pull --rm -t $(fullImageName) .

push: build
	@docker push $(fullImageName)

test: build
	@docker run --rm -d --name slack-notifier -p 8080:8080 -e NOTIFIER_WEBHOOKURL=$(webHookUrl) $(fullImageName)

minio: test
	@docker run -d --rm --privileged -p 9000:9000 --name minio -v $(PWD):/data -v $(PWD):/root/.minio --link slack-notifier:notifier minio/minio server /data