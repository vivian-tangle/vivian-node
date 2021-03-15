package network

import (
	"context"
	"log"
	"time"

	ggio "github.com/gogo/protobuf/io"
	proto "github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	p2p "github.com/vivian-tangle/vivian-client/network/pb"
)

// Node is the structure for lip-p2p node
type Node struct {
	host.Host     // lib-p2p host
	*PingProtocol // Ping protocol implementation
	Network       *Network
}

// AuthenticateMessage authenticate incoming p2p message
// message: a protobuf go data object
// data: common p2p message data
func (n *Node) AuthenticateMessage(message proto.Message, data *p2p.MessageData) bool {
	// store a temp ref to signature and remove it from message data
	// sign is a string to allow easy reset to zero-value (empty string)
	sign := data.Sign
	data.Sign = nil

	// marshall data without the signature to protobuf3 binary format
	bin, err := proto.Marshal(message)
	if err != nil {
		log.Println(err, "failed to marshal pb message")
		return false
	}

	// restore sig in message data (for possible future use)
	data.Sign = sign

	// restore peer id binary format from base58 encoded node id data
	peerID, err := peer.IDB58Decode(data.NodeID)
	if err != nil {
		log.Println(err, "failed to decode node id from base58")
		return false
	}

	// verify the data was authored by the signing peer identified by the public key
	// and signature included in the message
	return n.verifyData(bin, []byte(sign), peerID, data.NodePubKey)
}

// sign an outgoing p2p message payload
func (n *Node) signProtoMessage(message proto.Message) ([]byte, error) {
	data, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	return n.signData(data)
}

// sign a binary data using the local node's private key
func (n *Node) signData(data []byte) ([]byte, error) {
	key := n.Peerstore().PrivKey(n.ID())
	res, err := key.Sign(data)
	return res, err
}

// Verify incoming p2p message data integrity
// data: data to verify
// signature: author signature provided in the message payload
// peerID: author peer id from the message payload
// pubKeyData: author public key from the message payload
func (n *Node) verifyData(data []byte, signature []byte, peerID peer.ID, pubKeyData []byte) bool {
	key, err := crypto.UnmarshalPublicKey(pubKeyData)
	if err != nil {
		log.Println(err, "Failed to extract key from message key data")
		return false
	}

	// extract node id from the provided public key
	idFromKey, err := peer.IDFromPublicKey(key)
	if err != nil {
		log.Println(err, "Failed to extract peer id from public key")
		return false
	}

	// verify that message author node id matches the provided node public key
	if idFromKey != peerID {
		log.Println(err, "Node id and provided public key mismatch")
		return false
	}

	res, err := key.Verify(data, signature)
	if err != nil {
		log.Println(err, "Error authenticcating data")
		return false
	}

	return res
}

// NewMessageData is a helper method to generate message data shared between all node's p2p protocols
// messageID: unique for requests, copied from request for response
func (n *Node) NewMessageData(messageID string, gossip bool) *p2p.MessageData {
	// Add protobuf bin data for message author public key
	// this is useful for authenticating messages forwarded by a node authored by another node
	nodePubKey, err := n.Peerstore().PubKey(n.ID()).Bytes()
	if err != nil {
		panic("Failed to get public key for sender from local peer store.")
	}
	return &p2p.MessageData{
		ClientVersion: n.Network.Config.NodeClientVersion,
		NodeID:        peer.IDB58Encode(n.ID()),
		NodePubKey:    nodePubKey,
		Timestamp:     time.Now().Unix(),
		Id:            messageID,
		Gossip:        gossip,
	}

}

// sendProtoMessage is a helper method to write a protobuf go data object to a network stream
// data: reference of protobuf go data object to send (not the object itself)
// s: network stream to write the data to
func (n *Node) sendProtoMessage(id peer.ID, p protocol.ID, data proto.Message) bool {
	s, err := n.NewStream(context.Background(), id, p)
	if err != nil {
		log.Println(err)
	}
	defer s.Close()

	writer := ggio.NewFullWriter(s)
	err = writer.WriteMsg(data)
	if err != nil {
		log.Println(err)
		s.Reset()
		return false
	}

	return true
}
