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

## Running

```bash
./jellyfin-broadcaster -networks "192.168.88.255:7359" -interval 30s
```

Or build and run in one step:
```bash
go run . -networks "192.168.88.255:7359" -interval 30s
```

