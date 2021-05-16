#!/bin/bash

echo "Removing local terraform state"
rm -rf .terraform.lock.hcl
rm -rf .terraform
rm -rf terraform.tfstate
rm -rf terraform.tfstate.backup
