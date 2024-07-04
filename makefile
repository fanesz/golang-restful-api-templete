dev:
	cls
	go run main.go

build:
	go build -o ./dist/app

start-win:
	dist/app/backend.exe

start-linux:
	./dist/app

deploy:
	go build -o ./dist/app
	./dist/app