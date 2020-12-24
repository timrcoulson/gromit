provision:
	ssh pi@raspberrypi.lan 'curl -sL https://dtcooper.github.io/raspotify/install.sh | sh'
	scp raspotify pi@raspberrypi.lan:/home/pi/raspotify
	ssh pi@raspberrypi.lan 'sudo -u root mv /home/pi/raspotify /etc/default/raspotify'
	ssh pi@raspberrypi.lan 'sudo -u root systemctl restart raspotify'

run:
	export $(cat .env | xargs)
	go run .