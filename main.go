package main

import (
	"flag"
	"net"
	"time"
	"log"
	"strings"
)

const (
	DISCOVERY_MESSAGE = "Who is JellyfinServer?"
	DEFAULT_INTERVAL  = 30 * time.Second
)

var (
	remoteNetworks string
	interval       time.Duration
)

func init() {
	flag.StringVar(&remoteNetworks, "networks", "192.168.2.255:7359,192.168.3.255:7359", "Comma-separated list of remote network broadcast addresses and ports (e.g., 192.168.2.255:7359)")
	flag.DurationVar(&interval, "interval", DEFAULT_INTERVAL, "Interval between discovery packet sends (e.g., 30s, 1m)")
	flag.Parse()
}

func sendDiscoveryPacket(networks []string) error {
	// Get the local IP address
	localIP := getLocalIP()
	if localIP == nil {
		log.Printf("Error getting local IP: no valid IP found")
		return nil
	}

	// Create a UDP address for the local binding
	localAddr := &net.UDPAddr{
		IP:   localIP,
		Port: 7359,
	}

	// Listen on the local address
	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		log.Printf("Error listening on UDP: %v", err)
		return err
	}
	defer conn.Close()

	// Set write buffer and deadline
	conn.SetWriteBuffer(1024)
	err = conn.SetDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		log.Printf("Error setting deadline: %v", err)
		return err
	}

	for _, addr := range networks {
		log.Printf("Attempting to resolve %s", addr)
		remoteAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			log.Printf("Error resolving %s: %v", addr, err)
			continue
		}

		log.Printf("Resolved address: IP=%v, Port=%d", remoteAddr.IP, remoteAddr.Port)

		// Ensure the address is valid before sending
		if remoteAddr.IP == nil || remoteAddr.Port == 0 {
			log.Printf("Invalid address format for %s", addr)
			continue
		}

		log.Printf("Attempting to send to %s", addr)
		_, err = conn.WriteToUDP([]byte(DISCOVERY_MESSAGE), remoteAddr)
		if err != nil {
			log.Printf("Error sending to %s: %v", addr, err)
			continue
		}
		log.Printf("Sent discovery packet to %s", addr)
	}

	return nil
}

// getLocalIP returns the non-loopback local IP of the host
func getLocalIP() net.IP {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	for _, address := range addrs {
		// Check if the address is a valid unicast IP
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP
			}
		}
	}
	return nil
}

func main() {
	// Split the comma-separated list of networks
	networks := strings.Split(remoteNetworks, ",")
	if len(networks) == 0 {
		log.Fatal("No networks specified. Use -networks flag to provide a comma-separated list of addresses.")
	}

	// Trim any whitespace from the network addresses
	for i := range networks {
		networks[i] = strings.TrimSpace(networks[i])
	}

	for {
		log.Printf("Starting discovery cycle with networks: %v", networks)
		err := sendDiscoveryPacket(networks)
		if err != nil {
			log.Printf("Error in discovery packet send: %v", err)
		}
		time.Sleep(interval)
	}
}
