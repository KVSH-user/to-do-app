.PHONY:
.SILENT:

build:
	go build -o ./.bin/app cmd/todo/main.go
run: build
	./.bin/app