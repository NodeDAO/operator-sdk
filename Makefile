
.PHONY: go-bindata
go-bindata:
	# go install -a -v github.com/go-bindata/go-bindata/...@latest
	go generate ./config/var.go
