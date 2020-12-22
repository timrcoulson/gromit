.DEFAULT_GOAL := all

dev: | build-docker run-docker

build-docker:
	docker build -t gromit .
run-docker:
	docker run -it gromit

local:
	service cups start
	go run .

provision:
	ssh-copy-id pi@raspberrypi.lan
	scp provision.sh pi@raspberrypi.lan:/home/pi/provision.sh
	ssh pi@raspberrypi.lan 'sudo -u root /home/pi/provision.sh'

deploy:
	echo "Building binary"
	env GOOS=linux GOARCH=arm GOARM=5 go build

	echo "Stop service"
	ssh pi@raspberrypi.lan 'sudo -u root systemctl stop gromit'

	echo "Copying files"
	scp gromit pi@raspberrypi.lan:/home/pi/gromit
	scp gromit.service pi@raspberrypi.lan:/home/pi/gromit.service
	scp .env pi@raspberrypi.lan:/home/pi/.env

	echo "Starting service"
	ssh pi@raspberrypi.lan 'sudo -u root cp /home/pi/gromit.service /etc/systemd/system/gromit.service'
	ssh pi@raspberrypi.lan 'sudo -u root systemctl daemon-reload'
	ssh pi@raspberrypi.lan 'sudo -u root systemctl restart gromit'

logs:
	ssh pi@raspberrypi.lan 'sudo -u root journalctl -u gromit'