#!/usr/bin/env bash

apt update && apt install --no-install-recommends -y enscript cups cups-bsd ca-certificates bash jq && rm -rf /var/lib/apt/lists/*
service cups start

# Spotify
curl -sL https://dtcooper.github.io/raspotify/install.sh | sh

cat > /etc/default/raspotify <<- EndOfMessage
DEVICE_NAME="raspotify"
OPTIONS="--username tim.r.coulson --password VTvtQNyXpDbgwA8k9Etvd6Ke"
EndOfMessage

rm -rf /etc/gromit
mkdir /etc/gromit
chmod -R a+rwX /etc/gromit

systemctl restart raspotify
systemctl enable raspotify
