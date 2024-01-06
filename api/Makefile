.PHONY: build
build:
	CGO_ENABLED=0 go build -a -installsuffix cgo -o $(CURDIR)/bin $(CURDIR)/cmd/main.go

.PHONY: docker-build
docker-build:
	docker build --platform=linux/amd64  -t pavlokushneryk/crawler:latest .

.PHONY: docker-push
docker-push:
	docker push pavlokushneryk/crawler:latest