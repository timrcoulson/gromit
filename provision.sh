#!/usr/bin/env bash

apt update && apt install --no-install-recommends -y enscript cups cups-bsd ca-certificates bash jq && rm -rf /var/lib/apt/lists/*
service cups start
curl -sL https://dtcooper.github.io/raspotify/install.sh | sh
