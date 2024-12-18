PROJECT:=go-admin

.PHONY: build
build:
	goreleaser build --snapshot --clean

release:
	goreleaser release --snapshot --clean --skip publish
