# Hello World Example

## Quick start

### 1. Provider

The provider is not yet published so it must be built and installed locally. See [the provider readme](../../provider/README.md) for building and installing the provider.

### 2. Factorio Client + Server Setup

Run `./scripts/run.ps1`. It installs the mod to the current machine's Factorio client mods and sets up a headless factorio server, also running the mod.

To connect your client to the server choose "Multiplayer" > "Connect to address" > Use `127.0.0.1:34197` as the "IP address and port"

### 3. Terraform Run

[Get Terraform](https://www.terraform.io/downloads.html)

```
terraform init
terraform plan
terraform apply -auto-approve
```

### Cleanup

To wipe files created by the above operations you can use:

- `scripts/clean-tf.ps1`: Deletes just the terraform state.
- `scripts/clean-all.ps1`: Deletes the server state, and also removes the mod from your own Factorio client install.
