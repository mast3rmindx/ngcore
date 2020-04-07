package ngp2p

import (
	"io/ioutil"

	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/network"

	"github.com/ngchain/ngcore/ngp2p/pb"
)

// NotFound will reply NotFound message to remote node
func (w *Wired) NotFound(s network.Stream, uuid string) {
	log.Debugf("Sending notfound to %s. Message id: %s...", s.Conn().RemotePeer(), uuid)
	resp := &pb.Message{
		Header:  w.node.NewHeader(uuid),
		Payload: nil,
	}

	// sign the data
	signature, err := w.node.signMessage(resp)
	if err != nil {
		log.Errorf("failed to sign response")
		return
	}

	// add the signature to the message
	resp.Header.Sign = signature

	// send the response
	if ok := w.node.sendProtoMessage(s.Conn().RemotePeer(), notfoundMethod, resp); ok {
		log.Debugf("notfound to %s sent.", s.Conn().RemotePeer().String())
	}
}

// onNotFound is a remote notfound handler
func (w *Wired) onNotFound(s network.Stream) {
	buf, err := ioutil.ReadAll(s)
	if err != nil {
		_ = s.Reset()
		log.Error(err)
		return
	}

	// unmarshal it
	var data = &pb.Message{}
	err = proto.Unmarshal(buf, data)
	if err != nil {
		log.Error(err)
		return
	}

	if !w.node.verifyResponse(data) || !w.node.authenticateMessage(s.Conn().RemotePeer(), data) {
		log.Errorf("Failed to authenticate message")
		return
	}

	remoteID := s.Conn().RemotePeer()
	_ = s.Close()

	log.Debugf("Received notfound from %s. Message id:%s. Message: %s.", remoteID, data.Header.Uuid, data.Payload)
}
