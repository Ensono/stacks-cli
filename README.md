# Stacks CLI

The new home of the Stacks CLI for [Ensono Stacks](https://stacks.ensono.com).

A powerful project scaffolding tool that generates complete application templates with infrastructure, pipelines, and supporting files. Supports multiple frameworks (dotnet, java, nx) and cloud platforms (Azure, AWS).

## Installation

### Download Pre-built Binaries

Download the latest release from the [GitHub Releases page](https://github.com/Ensono/stacks-cli/releases).

#### Linux (AMD64)

```bash
curl -L -o stacks-cli https://github.com/Ensono/stacks-cli/releases/latest/download/stacks-cli-linux-amd64-$(curl -s https://api.github.com/repos/Ensono/stacks-cli/releases/latest | grep '"tag_name"' | cut -d '"' -f 4 | sed 's/^v//')
chmod +x stacks-cli
sudo mv stacks-cli /usr/local/bin/
stacks-cli version
```

#### macOS (Intel)

```bash
curl -L -o stacks-cli https://github.com/Ensono/stacks-cli/releases/latest/download/stacks-cli-darwin-amd64-$(curl -s https://api.github.com/repos/Ensono/stacks-cli/releases/latest | grep '"tag_name"' | cut -d '"' -f 4 | sed 's/^v//')
chmod +x stacks-cli
sudo mv stacks-cli /usr/local/bin/
stacks-cli version
```

#### macOS (Apple Silicon)

```bash
curl -L -o stacks-cli https://github.com/Ensono/stacks-cli/releases/latest/download/stacks-cli-darwin-arm64-$(curl -s https://api.github.com/repos/Ensono/stacks-cli/releases/latest | grep '"tag_name"' | cut -d '"' -f 4 | sed 's/^v//')
chmod +x stacks-cli
sudo mv stacks-cli /usr/local/bin/
stacks-cli version
```

#### Windows

**PowerShell (Recommended)**

```powershell
$version = (Invoke-RestMethod -Uri "https://api.github.com/repos/Ensono/stacks-cli/releases/latest").tag_name -replace '^v', ''
$url = "https://github.com/Ensono/stacks-cli/releases/latest/download/stacks-cli-windows-amd64-$version.exe"
Invoke-WebRequest -Uri $url -OutFile "stacks-cli.exe"
$userBin = "$env:USERPROFILE\bin"
if (!(Test-Path $userBin)) { New-Item -Path $userBin -ItemType Directory }
Move-Item "stacks-cli.exe" "$userBin\stacks-cli.exe"
$env:PATH = "$userBin;$env:PATH"
stacks-cli version
```

### Manual Installation

1. Download `stacks-cli-windows-amd64-{version}.exe` from the [releases page](https://github.com/Ensono/stacks-cli/releases/latest)
2. Rename the file to `stacks-cli.exe`
3. Add the executable to your PATH or place it in a directory that's already in your PATH
4. Open a new command prompt or PowerShell window and verify: `stacks-cli version`

### Manual Installation (Specific Version)

To install a specific version, replace `{version}` with the desired version number (e.g., `0.4.53`):

- **Linux**: `stacks-cli-linux-amd64-{version}`
- **macOS Intel**: `stacks-cli-darwin-amd64-{version}`
- **macOS Apple Silicon**: `stacks-cli-darwin-arm64-{version}`
- **Windows**: `stacks-cli-windows-amd64-{version}.exe`

## Quick Start

```bash
stacks-cli interactive
stacks-cli scaffold --name my-project --framework dotnet --cloud azure
stacks-cli --help
```

## Documentation

Documentation is stored with the code in Asciidoc format. It is in the `docs/` directory of the repository.

Each time a build is run a PDF file is generated as well as a set of Markdown files.

It is possible to run the documentation locally using Hugo in Docker. Due to the minimal support for Asciidoc in Hugo, a custom image has been built to run the website. Run the following command to run a local web server of the documentation.

```bash
docker run --rm -it -v ${PWD}/docs:/hugo-project/content/docs -v ${PWD}:/repo -p 1313:1313 russellseymour/hugo-docker
```
