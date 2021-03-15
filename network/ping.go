package network

import (
	"fmt"
	"io/ioutil"
	"log"

	proto "github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
	p2p "github.com/vivian-tangle/vivian-client/network/pb"
)

// PingProtocol type
type PingProtocol struct {
	Node         *Node                       // Local host
	Requests     map[string]*p2p.PingRequest // Used to access request data from response handlers
	PingRequest  protocol.ID                 // Ping request version
	PingResponse protocol.ID                 // Ping response version
	Done         chan bool                   // Only for demo purpose to stop main from terminating
}

// remote peer request handler
func (p *PingProtocol) onPingRequest(s network.Stream) {
	// get request data
	data := &p2p.PingRequest{}
	buf, err := ioutil.ReadAll(s)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}
	s.Close()

	// unmarshal it
	proto.Unmarshal(buf, data)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("%s: Sending ping response to %s/ Message id: %s...", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.MessageData.Id)

	valid := p.Node.AuthenticateMessage(data, data.MessageData)

	if !valid {
		log.Println("Failed to authenticate message")
		return
	}

	// generate response message
	log.Printf("%s: Sending ping response to %s. Message id: %s...", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.MessageData.Id)

	resp := &p2p.PingResponse{MessageData: p.Node.NewMessageData(data.MessageData.Id, false),
		Message: fmt.Sprintf("Ping response from %s", p.Node.ID())}

	// sign the data
	signature, err := p.Node.signProtoMessage(resp)
	if err != nil {
		log.Println("Failed to sign response")
		return
	}

	// add the signature to the message
	resp.MessageData.Sign = signature

	// send the response
	ok := p.Node.sendProtoMessage(s.Conn().RemotePeer(), p.PingResponse, resp)

	if ok {
		log.Printf("%s: Ping response to %s sent.", s.Conn().LocalPeer().String(), s.Conn().RemotePeer().String())
	}
}

// remote ping response handler
func (p *PingProtocol) onPingResponse(s network.Stream) {
	data := &p2p.PingRequest{}
	buf, err := ioutil.ReadAll(s)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}
	s.Close()

	// unmarshal it
	proto.Unmarshal(buf, data)
	if err != nil {
		log.Println(err)
		return
	}

	valid := p.Node.AuthenticateMessage(data, data.MessageData)

	if !valid {
		log.Println("Failed to authenticate message")
		return
	}

	// locate requests data and remove it if found
	_, ok := p.Requests[data.MessageData.Id]
	if ok {
		// remove requests from map as we have processed it here
		delete(p.Requests, data.MessageData.Id)
	} else {
		log.Println("Failed to locate request data object for response")
		return
	}

	log.Printf("%s: Receive ping response from %s. Message id:%s. Message: %s.", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.MessageData.Id, data.Message)
	p.Done <- true
}

// Ping is the function for sending the ping message
func (p *PingProtocol) Ping(host host.Host) bool {
	log.Printf("%s: Sending ping to: %s...", p.Node.ID(), host.ID())

	// create message data
	req := &p2p.PingRequest{MessageData: p.Node.NewMessageData(uuid.New().String(), false),
		Message: fmt.Sprintf("Ping from %s", p.Node.ID())}

	// sign the data
	signature, err := p.Node.signProtoMessage(req)
	if err != nil {
		log.Println("failed to sign pb data")
		return false
	}

	// add the signature to the message
	req.MessageData.Sign = signature

	ok := p.Node.sendProtoMessage(host.ID(), p.PingRequest, req)
	if !ok {
		return false
	}

	// store ref request so response handler has access to it
	p.Requests[req.MessageData.Id] = req
	log.Printf("%s: Ping to: %s was sent. Message Id: %s, Message: %s", p.Node.ID(), host.ID(), req.MessageData.Id, req.Message)
	return true
}
