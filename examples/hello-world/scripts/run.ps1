
& $PSScriptRoot\clean-all.ps1

# Copy the mod to your client mods folder
cp -r ../../mod/terraform-crud-api $env:APPDATA/Factorio/mods
# Create a folder to store the Factorio server data
mkdir factorio-volume
# Copy the factorio mod to the mods directory
mkdir factorio-volume/mods
cp -r ../../mod/terraform-crud-api factorio-volume/mods
# Configure the rcon pw
mkdir factorio-volume/config
Write-Output "SOMEPASSWORD" | Out-File -Encoding ASCII -NoNewLine factorio-volume/config/rconpw
# Run factorio server
docker run -it -p 127.0.0.1:34197:34197/udp -p 127.0.0.1:27015:27015/tcp -v "${PWD}/factorio-volume:/factorio" factoriotools/factorio:1.1.33