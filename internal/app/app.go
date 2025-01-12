package app

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"github.com/osamikoyo/chat-to-chat/internal/host"
	"github.com/osamikoyo/chat-to-chat/pkg/loger"
	"strings"

	"github.com/libp2p/go-libp2p/core/peer"
)

type App struct{
	Host *host.ChatHost
	Logger loger.Logger
}

func Init() *App {
	port := flag.Int("port", 0, "Port to listen on")
	flag.Parse()

	if *port == 0 {
		fmt.Println("Please specify a port number using -port flag")
		return &App{}
	}

	ctx := context.Background()

	// Create a new chat host
	chatHost, err := host.NewChatHost(ctx, *port)
	if err != nil {
		loger.New().Error().Err(err)
		return &App{}
	}
	defer chatHost.Close()

	return &App{
		Host: chatHost,
		Logger: loger.New(),
	}
}

func (a *App) Run() {
	go func() {
		for msg := range a.Host.GetMessages() {
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

		if strings.HasPrefix(input, "/connect ") {
			peerAddr := strings.TrimPrefix(input, "/connect ")
			err := a.Host.Connect(peerAddr)
			if err != nil {
				a.Logger.Info().Msg(fmt.Sprintf("Failed to connect: %v\n", err))
			} else {
				a.Logger.Info().Msg("Connected successfully!")
				// Extract peer ID from the multiaddress
				parts := strings.Split(peerAddr, "/p2p/")
				if len(parts) == 2 {
					currentPeer, err = peer.Decode(parts[1])
					if err != nil {
						a.Logger.Info().Msg(fmt.Sprintf("Failed to decode peer ID: %v\n", err))
					}
				}
			}
		} else if currentPeer != "" {
			// Send message to the connected peer
			err := a.Host.SendMessage(currentPeer, input)
			if err != nil {
				a.Logger.Info().Msg(fmt.Sprintf("Failed to send message: %v\n", err))
			}
		} else {
			a.Logger.Info().Msg("Not connected to any peer. Use /connect <peer-address> first")
		}

		fmt.Print("> ")
	}
}