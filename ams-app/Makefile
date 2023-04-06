IMAGE ?= quay.io/zroubalik/ams-app

.PHONY: image
image:
	docker build -f Dockerfile -t ${IMAGE}:latest .
	docker push ${IMAGE}:latest

.PHONY: dev
dev:
	docker build -f Dockerfile.dev -t ${IMAGE}:dev .
	docker push ${IMAGE}:dev