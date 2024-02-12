.PHONY: build run build-image start-container
.SILENT:

build:
	go build -o ./.bin/app cmd/todo/main.go

run: build
	./.bin/app

build-image:
	docker build -t kvshuser/todolist:v0.1 .

start-container:
	pwd
	ls -la
	docker run --env-file .env -p 8000:8000 kvshuser/todolist:v0.1
