# terraform-provider-factorio

This is the Terraform provider which talks to the Factorio API.

It is recommended to open vscode to this folder specifically to get the full devcontainer golang integration.

## Local Development

### Build

Change GOOS and GOARCH to match your setup:

```bash
# For Windows
GOOS=windows GOARCH=amd64 bash -c 'go build -o build/$GOOS/$GOARCH/terraform-provider-factorio.exe'
# Or for Mac
GOOS=darwin GOARCH=amd64 bash -c 'go build -o build/$GOOS/$GOARCH/terraform-provider-factorio'
```

### Install

Copy to your tf plugins dir. Paths vary depending on your OS/ARCH.

#### Windows

```powershell
mkdir $env:APPDATA\terraform.d\plugins\registry.terraform.io\efokschaner\factorio\0.1\windows_amd64
cp build/windows/amd64/terraform-provider-factorio.exe $env:APPDATA\terraform.d\plugins\registry.terraform.io\efokschaner\factorio\0.1\windows_amd64
```

#### macos

```bash
mkdir -p $HOME/.terraform.d/plugins/registry.terraform.io/efokschaner/factorio/0.1/darwin_amd64
cp build/darwin/amd64/terraform-provider-factorio $HOME/.terraform.d/plugins/registry.terraform.io/efokschaner/factorio/0.1/darwin_amd64
```

### tfplugindocs

Docs are generated by https://github.com/hashicorp/terraform-plugin-docs . However a couple of tricks are needed as we are using repo structure that is not precisely what's expected by tfplugindocs

- Hack `providerName` in `internal/provider/generate.go` so it's `"terraform-provider-factorio"`
- Re-run `go install` in `cmd/tfplugindocs` to re-install the patched tool
- Copy the docs folder it generates to the root of this repo
