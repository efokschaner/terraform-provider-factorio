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
resource "factorio_hello" "a-greeting" {
  create_as_ghost = false
}