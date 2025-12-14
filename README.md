# Jellyfin Broadcaster

This repository contains a Go script designed to broadcast Jellyfin server discovery packets to remote networks, facilitating client discovery across different subnets or houses. The script is particularly useful when clients are not directly connected to the Tailscale network but need to discover the Jellyfin server.

## Overview

Jellyfin uses UDP port 7359 for server discovery, sending broadcast packets with the message "Who is JellyfinServer?" to facilitate client connections. This script replicates that behavior, sending the discovery packets from a Tailscale box to remote networks specified via command-line arguments.

## Features

- Sends Jellyfin discovery packets to multiple remote networks.
- Configurable via command-line arguments for flexibility.
- Lightweight and resource-efficient, written in Go.
- Periodic broadcasting with adjustable intervals.

## Prerequisites

- Go (Golang) installed on your system.
- A Tailscale box or similar device capable of sending UDP packets.
- Network access to the remote networks (e.g., 192.168.88.255).

## Installation

1. Clone this repository:
   ```bash
   git clone https://github.com/your-username/jellyfin-broadcaster.git
   cd jellyfin-broadcaster
   go build
   ```

## Running

### Manual Execution

Run directly:
```bash
./jellyfin-broadcaster -networks "192.168.88.255:7359" -interval 30s
```

Or build and run in one step:
```bash
go run . -networks "192.168.88.255:7359" -interval 30s
```

### Systemd Service

To run as a systemd service:

1. **Build and install the binary:**
   ```bash
   go build
   sudo cp jellyfin-broadcaster /usr/local/bin/
   sudo chmod +x /usr/local/bin/jellyfin-broadcaster
   ```

2. **Edit the systemd unit file** (`jellyfin-broadcaster.service`) to customize:
   - The `-networks` parameter with your target broadcast addresses
   - The `-interval` parameter for your desired broadcast frequency
   - Optionally uncomment and set the `User` field to run as a specific user (requires network permissions)

3. **Install and enable the service:**
   ```bash
   sudo cp jellyfin-broadcaster.service /etc/systemd/system/
   sudo systemctl daemon-reload
   sudo systemctl enable jellyfin-broadcaster.service
   sudo systemctl start jellyfin-broadcaster.service
   ```

4. **Check service status:**
   ```bash
   sudo systemctl status jellyfin-broadcaster.service
   ```

5. **View logs:**
   ```bash
   sudo journalctl -u jellyfin-broadcaster.service -f
   ```

**Service Management:**
- Stop: `sudo systemctl stop jellyfin-broadcaster.service`
- Restart: `sudo systemctl restart jellyfin-broadcaster.service`
- Disable: `sudo systemctl disable jellyfin-broadcaster.service`

