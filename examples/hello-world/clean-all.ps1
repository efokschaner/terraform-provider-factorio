& .\clean-terraform.ps1

write-host "Removing local client mod"
remove-item -force -recurse -ErrorAction Ignore $env:APPDATA/Factorio/mods/terraform-crud-api

write-host "Removing server persistent volume"
remove-item -force -recurse -ErrorAction Ignore factorio-volume