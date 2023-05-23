VERSION=latest
USER=nodedao

.PHONY: build-local
build-local:
	CGO_ENABLED=0  GOOS=linux  GOARCH=amd64 go build -o ./build/operator-sdk main.go
	docker buildx build --platform=linux/amd64 -t $(USER)/operator-sdk:$(VERSION) -f Dockerfile-local-example ./build

.PHONY: build
build:
	docker buildx build --platform=linux/amd64 -t $(USER)/operator-sdk:$(VERSION)  .

.PHONY: push
push:
	docker push $(USER)/operator-sdk:$(VERSION)

.PHONY: go-bindata
go-bindata:
	# go install -a -v github.com/go-bindata/go-bindata/...@latest
	go generate ./config/var.go