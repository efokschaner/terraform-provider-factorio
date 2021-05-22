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

# Example of state fetching
data "factorio_players" "all" {}

output "all_players" {
  value = data.factorio_players.all.players
}

# Example of resource creating
resource "factorio_entity" "a-furnace" {
  surface = "nauvis"
  name = "stone-furnace"
  position {
    x = 1
    y = 2
  }
  direction = "east"
  force = "player"
}
