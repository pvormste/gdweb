# gdweb

Zero-config HTTPS server for Godot 4 web exports. Serves your game over TLS so you can test it on other devices on your local network (Godot requires a secure context for network testing).

## Install

```bash
go install github.com/pvormste/gdweb@latest
```

## Usage

From your Godot web export directory:

```bash
gdweb
```

This serves the current directory over HTTPS on port 8443. The server prints local and network URLs so you can open the game on your phone, tablet, or another computer.

### QR codes

Visit `/qr` (e.g. `https://localhost:8443/qr`) for a page with tabbed QR codes for each address. Scan the appropriate QR code with your phone to open the game without typing the URL.

### Options

| Flag   | Default | Description                    |
|--------|---------|--------------------------------|
| `-port`| 8443    | HTTPS port to listen on       |
| `-dir` | .       | Directory to serve            |
| `-host`| 0.0.0.0 | Address to bind to            |
| `-open`| false   | Open browser on startup       |

### Example

```bash
# Serve current directory on default port
gdweb

# Serve a specific export folder on port 9000
gdweb -dir ./exports/web -port 9000

# Open browser automatically
gdweb -open
```

## Requirements

- Godot 4 web export (single-threaded or multi-threaded)
- Self-signed certificate: accept the browser warning once to proceed
