write-host "Removing local terraform state"
remove-item -force -recurse -ErrorAction Ignore .terraform.lock.hcl
remove-item -force -recurse -ErrorAction Ignore .terraform
remove-item -force -recurse -ErrorAction Ignore terraform.tfstate
remove-item -force -recurse -ErrorAction Ignore terraform.tfstate.backup
