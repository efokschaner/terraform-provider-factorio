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

# Example of resource creating
resource "factorio_entity" "a-furnace" {
  surface = "nauvis"
  name = "stone-furnace"
  position {
    x = 1
    y = 2
  }
  direction = "north"
  force = "player"
}

# Creating a ghost, requires entity_specific_parameters
resource "factorio_entity" "a-ghost-furnace" {
  surface = "nauvis"
  name = "entity-ghost"
  position {
    x = 3
    y = 2
  }
  direction = "north"
  force = "player"
  entity_specific_parameters = {
    inner_name = "stone-furnace"
    expires = false
  }
}

// Example of using a Factorio infrastructure module
module "hello" {
  count = length(data.factorio_players.all.players)
  source = "./modules/hello"
  position = {
    x = data.factorio_players.all.players[count.index].position.0.x
    y = data.factorio_players.all.players[count.index].position.0.y
  }
  direction = "east"
}