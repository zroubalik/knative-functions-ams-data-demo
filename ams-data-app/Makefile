IMAGE ?= quay.io/zroubalik/ams-app

image: backend frontend

.PHONY: backend
backend:
	docker build -f Dockerfile.backend -t ${IMAGE}:backend .
	docker push ${IMAGE}:backend

.PHONY: frontend
frontend:
	docker build -f Dockerfile.frontend -t ${IMAGE}:frontend .
	docker push ${IMAGE}:frontend

run:
	npm start

run-images:
	docker run -d -p 3333:3333 ${IMAGE}:backend   
	docker run -d -p 3000:3000 ${IMAGE}:frontend 


.PHONY: fff
fff:
#	docker build -f Dockerfile -t ${IMAGE}:frontendx .
	docker push ${IMAGE}:frontendx
