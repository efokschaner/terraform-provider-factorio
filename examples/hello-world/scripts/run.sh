#!/bin/bash

SCRIPTS_DIR=`dirname $0`

$SCRIPTS_DIR/clean-all.sh

# Copy the mod to your client mods folder
cp -r ../../mod/terraform-crud-api "$HOME/Library/Application Support/factorio/mods"
# Create a folder to store the Factorio server data
mkdir factorio-volume
# Copy the factorio mod to the mods directory
mkdir factorio-volume/mods
cp -r ../../mod/terraform-crud-api factorio-volume/mods
# Configure the rcon pw
mkdir factorio-volume/config
echo "SOMEPASSWORD" > factorio-volume/config/rconpw
# Run factorio server
docker run -it -p 127.0.0.1:34197:34197/udp -p 127.0.0.1:27015:27015/tcp -v "$(pwd)/factorio-volume:/factorio" factoriotools/factorio:1.1.33