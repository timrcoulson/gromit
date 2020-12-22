.DEFAULT_GOAL := all

all: | build-docker run-docker

build-docker:
	docker build -t gromit .
run-docker:
	docker run --net=host -it gromit
run:
	service cups start
	go run ./main.go
