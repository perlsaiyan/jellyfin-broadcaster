package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	PortNumber        = 7359
	DiscoveryMessage  = "who is JellyfinServer?"
)

// ServerDiscoveryInfo matches the C# ServerDiscoveryInfo model
type ServerDiscoveryInfo struct {
	Address        string  `json:"Address"`
	Id             string  `json:"Id"`
	Name           string  `json:"Name"`
	EndpointAddress *string `json:"EndpointAddress,omitempty"`
}

var (
	serverURL        string
	serverID         string
	serverName       string
	endpointAddress  string
)

func init() {
	flag.StringVar(&serverURL, "server-url", "", "The Jellyfin server URL (e.g., http://192.168.1.100:8096) (required)")
	flag.StringVar(&serverID, "server-id", "", "Server identifier (GUID/UUID format) (required)")
	flag.StringVar(&serverName, "server-name", "", "Friendly server name (required)")
	flag.StringVar(&endpointAddress, "endpoint-address", "", "Optional endpoint address")
	flag.Parse()
}

func main() {
	// Validate required flags
	if serverURL == "" {
		log.Fatal("Error: -server-url is required")
	}
	if serverID == "" {
		log.Fatal("Error: -server-id is required")
	}
	if serverName == "" {
		log.Fatal("Error: -server-name is required")
	}

	// Create context that cancels on interrupt
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down...")
		cancel()
	}()

	// Start listening for discovery messages
	if err := listenForAutoDiscoveryMessage(ctx); err != nil {
		log.Fatalf("Failed to start discovery server: %v", err)
	}
}

func listenForAutoDiscoveryMessage(ctx context.Context) error {
	// Listen on all interfaces (0.0.0.0) on port 7359
	listenAddr := &net.UDPAddr{
		IP:   net.IPv4zero, // 0.0.0.0
		Port: PortNumber,
	}

	conn, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Printf("Listening for discovery requests on %s:%d", listenAddr.IP, PortNumber)

	buffer := make([]byte, 1024)

	for {
		select {
		case <-ctx.Done():
			log.Println("Discovery socket operation cancelled")
			return nil
		default:
			// Set read deadline to allow periodic context checks
			// This allows us to check ctx.Done() periodically
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))

			n, remoteAddr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				// Check if it's a timeout (expected for periodic context checks)
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				// Log other errors but continue listening
				log.Printf("Failed to receive data from socket: %v", err)
				continue
			}

			// Check if the message contains the discovery request (case-insensitive)
			message := strings.ToLower(string(buffer[:n]))
			if strings.Contains(message, strings.ToLower(DiscoveryMessage)) {
				log.Printf("Received discovery request from %s", remoteAddr)
				if err := respondToDiscoveryMessage(remoteAddr, conn); err != nil {
					log.Printf("Error sending response message: %v", err)
				}
			}
		}
	}
}

func respondToDiscoveryMessage(endpoint *net.UDPAddr, conn *net.UDPConn) error {
	// Create the response with server information
	var endpointAddrPtr *string
	if endpointAddress != "" {
		endpointAddrPtr = &endpointAddress
	}

	response := ServerDiscoveryInfo{
		Address:        serverURL,
		Id:             serverID,
		Name:           serverName,
		EndpointAddress: endpointAddrPtr,
	}

	// Serialize to JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return err
	}

	log.Printf("Sending AutoDiscovery response to %s", endpoint)
	_, err = conn.WriteToUDP(responseJSON, endpoint)
	return err
}
