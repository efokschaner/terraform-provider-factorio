# Terraform Provider for Factorio

"Infrastructure as Code" for your factory.

https://user-images.githubusercontent.com/1409112/119280384-0a067680-bbe6-11eb-8610-10a3f5a9eeb5.mp4

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
