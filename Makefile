REPO ?= ntate22/logging-generator
DEV_TAG ?= dev
COMMIT ?= $(shell git rev-parse --short HEAD)
SHELL=/bin/bash -o pipefail

image: dev-image
	docker build -t $(REPO):$(COMMIT) --build-arg=DEV_IMAGE=$(REPO):$(DEV_TAG) .
ifdef latest
	docker tag $(REPO):$(COMMIT) $(REPO):latest
endif

dev-image:
	docker build -t $(REPO):$(DEV_TAG) -f Dockerfile.dev .

image-push: image
	docker push $(REPO):$(COMMIT)
ifdef latest
	docker push $(REPO):latest
endif

test: dev-image
	docker run $(REPO):$(DEV_TAG) go test github.com/nicktate/logging-generator/...

.PHONY: dev-image image
