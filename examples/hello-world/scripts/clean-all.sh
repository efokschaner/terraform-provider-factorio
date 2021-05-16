#!/bin/bash

SCRIPTS_DIR=`dirname $0`

$SCRIPTS_DIR/clean-tf.sh

echo "Removing local client mod"
rm -rf "$HOME/Library/Application Support/factorio/mods/terraform-crud-api"

echo "Removing server persistent volume"
rm -rf factorio-volume