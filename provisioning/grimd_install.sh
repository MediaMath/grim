#!/bin/bash

set -eu
export DEBIAN_FRONTEND=noninteractive

### install prerequisites ###
sudo apt-get update
sudo apt-get install -y -q zip

### create grim user ###
useradd -s /bin/bash -m -d /var/lib/grim grim

### create expected directories ###
mkdir -p /opt/grimd # install dir
mkdir -p /var/log/grim # logs
mkdir -p /var/tmp/grim # build folders
mkdir -p /etc/grim # config dir

### install grim ###
unzip -d/opt/grimd /tmp/grimd.zip

### write the service config ###
cat > /etc/systemd/system/grimd.service <<End-of-grimd.service
[Unit]
Description=Grimd

[Service]
Type=simple
ExecStart=/opt/grimd/grimd
User=grim

[Install]
WantedBy=multi-user.target
End-of-grimd.service

### enable grimd ###
systemctl enable grimd.service
