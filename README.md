# Terraform Provider for Factorio

"Infrastructure as Code" for your factory.

TODO: Add gifs.

_Current Status:_ Barely functional and mostly useless.

Inspired by the likes of:

- https://github.com/abesto/codetorio
- https://github.com/Redcrafter/verilog2factorio/

Only works with factorio multiplayer server, as it depends on remote control via RCON.
See [./examples/hello-world](./examples/hello-world) for more information on how to use.

## Repository Overview

- [`examples`](./examples): Examples using the provider
- [`mod`](./mod): The mod for factorio which provides an API for the provider.
- [`provider`](./provider): The Terraform provider.
