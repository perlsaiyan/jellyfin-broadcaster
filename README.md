# Jellyfin Broadcaster

This repository contains a Go server that responds to Jellyfin server discovery requests on UDP port 7359, matching the behavior of Jellyfin's `AutoDiscoveryHost` implementation. The server listens for discovery messages and responds with server information to facilitate client connections.

## Overview

Jellyfin uses UDP port 7359 for server discovery. Clients send broadcast packets with the message "who is JellyfinServer?" and servers respond with JSON containing server information (address, ID, name). This server implements the server-side discovery protocol, allowing Jellyfin clients to discover your server.

## Features

- Listens on UDP port 7359 for discovery requests
- Responds to discovery messages with JSON server information
- Configurable server URL, ID, and name via command-line flags
- Lightweight and resource-efficient, written in Go
- Matches Jellyfin's AutoDiscoveryHost implementation

## Prerequisites

- Go (Golang) installed on your system
- Network access to bind to UDP port 7359

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
./jellyfin-broadcaster -server-url "http://192.168.1.100:8096" -server-id "your-server-id" -server-name "My Jellyfin Server"
```

Or build and run in one step:
```bash
go run . -server-url "http://192.168.1.100:8096" -server-id "your-server-id" -server-name "My Jellyfin Server"
```

With optional endpoint address:
```bash
./jellyfin-broadcaster -server-url "http://192.168.1.100:8096" -server-id "your-server-id" -server-name "My Jellyfin Server" -endpoint-address "192.168.1.100"
```

### Command-Line Flags

- `-server-url` (required): The Jellyfin server URL (e.g., `http://192.168.1.100:8096`)
- `-server-id` (required): Server identifier (GUID/UUID format)
- `-server-name` (required): Friendly server name
- `-endpoint-address` (optional): Endpoint address

### Systemd Service

To run as a systemd service:

1. **Build and install the binary:**
   ```bash
   go build
   sudo cp jellyfin-broadcaster /usr/local/bin/
   sudo chmod +x /usr/local/bin/jellyfin-broadcaster
   ```

2. **Edit the systemd unit file** (`jellyfin-broadcaster.service`) to customize:
   - The `-server-url` parameter with your Jellyfin server URL
   - The `-server-id` parameter with your server's unique identifier
   - The `-server-name` parameter with your desired server name
   - Optionally add `-endpoint-address` if needed
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

## How It Works

The server listens on UDP port 7359 (all interfaces, 0.0.0.0) for incoming discovery requests. When a client sends a message containing "who is JellyfinServer?" (case-insensitive), the server responds with a JSON payload containing:

- `Address`: The server URL
- `Id`: The server identifier
- `Name`: The friendly server name
- `EndpointAddress`: Optional endpoint address

This matches the behavior of Jellyfin's `AutoDiscoveryHost.cs` implementation, allowing standard Jellyfin clients to discover your server.
