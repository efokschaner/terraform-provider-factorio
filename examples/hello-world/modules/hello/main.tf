terraform {
  required_providers {
    factorio = {
      version = "~> 0.1"
      source  = "efokschaner/factorio"
    }
  }
}

locals {
  // Offsets place the text with the center
  offset_x = -12
  offset_y = -10
  text = [
    [
      "X  X  XXX  X    X     XX ",
      "X  X  X    X    X    X  X",
      "XXXX  XXX  X    X    X  X",
      "X  X  X    X    X    X  X",
      "X  X  XXX  XXX  XXX   XX ",
    ],
    [
      "XXX  XXX    XX   X   X",
      "X    X  X  X  X  XX XX",
      "XXX  XXX   X  X  X X X",
      "X    X X   X  X  X   X",
      "X    X  X   XX   X   X",
    ],
    [
      "XXX  XXX  XXX   XXX    XX   XXX   XX   XXX   X   X",
      " X   X    X  X  X  X  X  X  X    X  X  X  X  XX XX",
      " X   XXX  XXX   XXX   XXXX  XXX  X  X  XXX   X X X",
      " X   X    X X   X X   X  X  X    X  X  X X   X   X",
      " X   XXX  X  X  X  X  X  X  X     XX   X  X  X   X",
    ],
  ]
  
  flat_text = flatten([
    for text_line in local.text:
      # Add 3 empty lines between each text line
      concat(text_line, ["", "", ""])
  ])

  pixels = flatten([
    for pixel_line_index, pixel_line in local.flat_text: [
      for pixel_index, pixel in regexall(".", pixel_line):
      {
        x = pixel_index
        y = pixel_line_index
      } if pixel == "X"
    ]
  ])

  belt_map = { for pixel in local.pixels :
    "_${pixel.x}_${pixel.y}" => pixel
  }
}

resource "factorio_entity" "belt" {
  for_each = local.belt_map
  surface = "nauvis"
  name = "transport-belt"
  position {
    x = each.value.x + local.offset_x + var.position.x
    y = each.value.y + local.offset_y + var.position.y
  }
  direction = var.direction
  force = "player"
}