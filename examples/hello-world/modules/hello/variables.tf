variable "position" {
  type = object({
    x = number
    y = number
  })
  default = {
    x = 0
    y = 0
  }
  description = "The location in the world at which the message is rendered"
}

variable "direction" {
  type = string
  default = "north"
  description = "The direction the conveyor belts are pointing"
}