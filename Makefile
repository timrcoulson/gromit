.DEFAULT_GOAL := all

all: | build-docker run-docker

deps:
	brew tap heroku/brew && brew install heroku
build-docker:
	docker build -t gromit .
run-docker:
	docker run -it gromit
run:
	cupsctl --debug-logging
	service cups start
	go run ./main.go
deploy:
	git push heroku main
