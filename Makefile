.DEFAULT_GOAL := all

dev: | build-docker run-docker

build-docker:
	docker build -t gromit .
run-docker:
	docker run -it gromit

local:
	service cups start
	HOST=http://localhost go run .

provision:
	ssh-copy-id pi@raspberrypi.local
	scp provision.sh pi@raspberrypi.local:/home/pi/provision.sh

	ssh pi@raspberrypi.local 'sudo -u root  /home/pi/provision.sh'

print:
	curl http://raspberrypi.local/print	

deploy:
	echo "Building binary"
	env GOOS=linux GOARCH=arm GOARM=5 go build

	echo "Stop service"
	ssh pi@raspberrypi.local 'sudo -u root systemctl stop gromit' || true

	echo "Copying files"
	scp gromit pi@raspberrypi.local:/home/pi/gromit
	scp conf/gromit.service pi@raspberrypi.local:/home/pi/gromit.service
	scp .env pi@raspberrypi.local:/home/pi/.env

	rsync -a --ignore-existing data pi@raspberrypi.local:/home/pi/

	rm gromit

	echo "Starting service"
	ssh pi@raspberrypi.local 'sudo -u root cp /home/pi/gromit.service /etc/systemd/system/gromit.service'
	ssh pi@raspberrypi.local 'sudo -u root systemctl daemon-reload'
	ssh pi@raspberrypi.local 'sudo -u root systemctl restart gromit'
	ssh pi@raspberrypi.local 'sudo -u root systemctl enable gromit'

logs:
	ssh pi@raspberrypi.local 'sudo -u root journalctl -u gromit'
