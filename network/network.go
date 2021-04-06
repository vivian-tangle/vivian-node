package network

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/protocol"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/vivian-tangle/vivian-client/config"
	p2p "github.com/vivian-tangle/vivian-client/network/pb"
)

// Network is the structure for storing the information of Libp2p node
type Network struct {
	Config *config.Config
}

// MakeRandomNode creates a lib-p2p host to listen to a port
func (nt *Network) MakeRandomNode(port int, done chan bool) *Node {
	// Ignoring most errors for brevity
	// See echo example for more details and better implementation
	priv, _, _ := crypto.GenerateKeyPair(crypto.Secp256k1, 256)
	listen, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port))
	host, _ := libp2p.New(
		context.Background(),
		libp2p.ListenAddrs(listen),
		libp2p.Identity(priv),
	)

	return nt.NewNode(host, done)
}

// NewNode creates a new node with its implemented protocols
func (nt *Network) NewNode(host host.Host, done chan bool) *Node {
	node := &Node{Host: host, Network: nt}
	node.PingProtocol = nt.NewPingProtocol(node, done)
	return node
}

// NewPingProtocol creates a new Ping protocol
func (nt *Network) NewPingProtocol(node *Node, done chan bool) *PingProtocol {
	p := &PingProtocol{
		Node:         node,
		Requests:     make(map[string]*p2p.PingRequest),
		PingRequest:  protocol.ID(nt.Config.PingRequestVersion),
		PingResponse: protocol.ID(nt.Config.PingResponseVersion),
		Done:         done,
	}
	node.SetStreamHandler(protocol.ID(nt.Config.PingRequestVersion), p.onPingRequest)
	node.SetStreamHandler(protocol.ID(nt.Config.PingResponseVersion), p.onPingResponse)

	return p
}
