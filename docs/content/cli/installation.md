---
title: Install MiniPaaS CLI
weight: 2
---

## Installation Options

You can install **MiniPaaS CLI** using one of the following methods:

---

### Option 1: Install via Go

If you have Go 1.25+ installed, use:

```bash
go install github.com/yunier-rojas/minipaas/minipaas-cli/cmd/minipaas@main
```

> This places the `minipaas` binary in `$(go env GOPATH)/bin`.
> Make sure that directory is in your system `PATH`:

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
```

---

### Option 2: Download Prebuilt Binaries

1. Go to the [GitHub Releases](https://github.com/yunier-rojas/minipaas/releases).

2. Download the binary for your platform:

	* `minipaas-linux-amd64`
	* `minipaas-darwin-arm64`
	* `minipaas-windows-amd64.exe`

3. Make it executable (Linux/macOS):

```bash
chmod +x minipaas-*
mv minipaas-* /usr/local/bin/minipaas
```

4. Confirm installation:

```bash
minipaas --help
```

---

### Option 3: Build from Source

If you prefer to build manually:

```bash
git clone https://github.com/yunier-rojas/minipaas.git
cd minipaas/minipaas-cli/cmd/minipaas
go mod tidy
go build
```

The command writes the `minipaas` binary to the current directory.

---

## Requirements

* Go **1.25+**
* Docker CLI
* An SSH client and key when targeting a remote host
* Unix-like shell or terminal (Linux, macOS, or Windows PowerShell)

---

Need help? Open an [issue](https://github.com/yunier-rojas/minipaas/issues) or [contact me](/contact/).
