module terraform-provider-factorio

go 1.16

require (
	github.com/gtaylor/factorio-rcon v0.0.0-20170109054031-61bdfe779ea6
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.6.1
)

replace github.com/gtaylor/factorio-rcon v0.0.0-20170109054031-61bdfe779ea6 => github.com/efokschaner/factorio-rcon v0.0.0-20210507061126-b1135a35d951
