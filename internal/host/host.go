package host

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// ChatHost represents a P2P chat host
type ChatHost struct {
	Host    host.Host
	ctx     context.Context
	msgChan chan ChatMessage
}

// ChatMessage represents a chat message with sender information
type ChatMessage struct {
	SenderID string
	Content  string
}

// NewChatHost creates a new P2P chat host
func NewChatHost(ctx context.Context, port int) (*ChatHost, error) {
	// Create a new libp2p host
	h, err := libp2p.New(
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port)),
	)
	if err != nil {
		return nil, err
	}

	chatHost := &ChatHost{
		Host:    h,
		ctx:     ctx,
		msgChan: make(chan ChatMessage, 100),
	}

	// Set stream handler for incoming chat messages
	h.SetStreamHandler("/chat/1.0.0", chatHost.handleStream)

	// Print host details
	addrs := h.Addrs()
	for _, addr := range addrs {
		fmt.Printf("P2P Address: %s/p2p/%s\n", addr, h.ID().String())
	}

	return chatHost, nil
}

// Connect connects to a peer using their multiaddress
func (ch *ChatHost) Connect(peerAddr string) error {
	maddr, err := multiaddr.NewMultiaddr(peerAddr)
	if err != nil {
		return fmt.Errorf("invalid peer address: %v", err)
	}

	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return fmt.Errorf("failed to get peer info: %v", err)
	}

	if err := ch.Host.Connect(ch.ctx, *info); err != nil {
		return fmt.Errorf("failed to connect to peer: %v", err)
	}

	return nil
}

// handleStream handles incoming message streams
func (ch *ChatHost) handleStream(s network.Stream) {
	go ch.readData(s)
}

// readData reads incoming messages from the stream
func (ch *ChatHost) readData(s network.Stream) {
	buf := make([]byte, 1024)
	for {
		n, err := s.Read(buf)
		if err != nil {
			s.Close()
			return
		}

		msg := ChatMessage{
			SenderID: s.Conn().RemotePeer().String(),
			Content:  string(buf[:n]),
		}
		ch.msgChan <- msg
	}
}

// SendMessage sends a message to a specific peer
func (ch *ChatHost) SendMessage(peerId peer.ID, message string) error {
	stream, err := ch.Host.NewStream(ch.ctx, peerId, "/chat/1.0.0")
	if err != nil {
		return fmt.Errorf("failed to create stream: %v", err)
	}
	defer stream.Close()

	_, err = stream.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	return nil
}

// GetMessages returns the channel for receiving messages
func (ch *ChatHost) GetMessages() <-chan ChatMessage {
	return ch.msgChan
}

// GetPeerID returns the host's peer ID
func (ch *ChatHost) GetPeerID() peer.ID {
	return ch.Host.ID()
}

// Close closes the host and all its connections
func (ch *ChatHost) Close() error {
	return ch.Host.Close()
}

