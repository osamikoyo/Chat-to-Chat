package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/osamikoyo/chat-to-chat/internal/data"
	"os"
	"strings"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/osamikoyo/chat-to-chat/internal/host"
)

func main() {
	storage, err := data.New()
	if err != nil{
		fmt.Print(err)
	}
	port := flag.Int("port", 0, "Port to listen on")
	flag.Parse()

	if *port == 0 {
		fmt.Println("Please specify a port number using -port flag")
		return
	}

	ctx := context.Background()

	// Create a new chat host
	chatHost, err := host.NewChatHost(ctx, *port)
	if err != nil {
		fmt.Printf("Failed to create chat host: %v\n", err)
		return
	}
	defer chatHost.Close()

	// Start message receiver
	go func() {
		for msg := range chatHost.GetMessages() {
			fmt.Printf("\n[%s]: %s\n> ", msg.SenderID, msg.Content)
		}
	}()

	// Main input loop
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Commands:")
	fmt.Println("/connect <peer-address> - Connect to a peer")
	fmt.Println("/quit - Exit the chat")
	fmt.Print("> ")

	var currentPeer peer.ID
	for scanner.Scan() {
		input := scanner.Text()

		if input == "/quit" {
			break
		}

		if strings.HasPrefix(input, "/history ") {
			peerAddr := strings.TrimPrefix(input, "/history ")
			addr := chatHost.Host.Addrs()[0]
			messages, err := storage.Get(10, peerAddr, fmt.Sprintf("%s/p2p/%s\n", addr, chatHost.Host.ID().String()))
			if err != nil{
				fmt.Print(err)
			}

			for _,  msg := range messages{
				if msg.Receiver == peerAddr {
					fmt.Printf("%s", msg.Content)
				} else {
					fmt.Printf("										%s", msg.Content)
				}
			}
		}

		if strings.HasPrefix(input, "/connect ") {
			peerAddr := strings.TrimPrefix(input, "/connect ")
			err := chatHost.Connect(peerAddr)
			if err != nil {
				fmt.Printf("Failed to connect: %v\n", err)
			} else {
				fmt.Println("Connected successfully!")
				// Extract peer ID from the multiaddress
				parts := strings.Split(peerAddr, "/p2p/")
				if len(parts) == 2 {
					currentPeer, err = peer.Decode(parts[1])
					if err != nil {
						fmt.Printf("Failed to decode peer ID: %v\n", err)
					}
				}
			}
		} else if currentPeer != "" {
			// Send message to the connected peer
			err := chatHost.SendMessage(currentPeer, input)
			if err != nil {
				fmt.Printf("Failed to send message: %v\n", err)
			}
		} else {
			fmt.Println("Not connected to any peer. Use /connect <peer-address> first")
		}

		fmt.Print("> ")
	}
}