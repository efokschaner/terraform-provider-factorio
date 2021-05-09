terraform {
  required_providers {
    factorio = {
      version = "~> 0.1"
      source  = "efokschaner/factorio"
    }
  }
}

provider "factorio" {
  rcon_host = "127.0.0.1:27015"
  rcon_pw  = "SOMEPASSWORD"
}


data "factorio_players" "all" {}

# Returns all players
output "all_players" {
  value = data.factorio_players.all.players
}