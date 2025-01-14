package consensus

import (
	"fmt"

	"github.com/ngchain/ngcore/ngp2p/message"
	"github.com/ngchain/ngcore/ngp2p/wired"

	core "github.com/libp2p/go-libp2p-core"

	"github.com/ngchain/ngcore/ngtypes"
)

// GetRemoteStatus just get the remote status from remote, and then put it into sync.store
func (mod *syncModule) getRemoteStatus(peerID core.PeerID) error {
	origin := mod.pow.Chain.GetOriginBlock()
	latest := mod.pow.Chain.GetLatestBlock()
	cp := mod.pow.Chain.GetLatestCheckpoint()

	id, stream := mod.localNode.SendPing(peerID, origin.GetHeight(), latest.GetHeight(), cp.Hash(), cp.GetActualDiff().Bytes())
	if stream == nil {
		log.Infof("failed to send ping, cannot get remote status from %s", peerID) // level down this
		return nil
	}

	reply, err := wired.ReceiveReply(id, stream)
	if err != nil {
		return err
	}

	switch reply.Header.MessageType {
	case message.MessageType_PONG:
		pongPayload, err := wired.DecodePongPayload(reply.Payload)
		if err != nil {
			return err
		}

		if _, exists := mod.store[peerID]; !exists {
			mod.putRemote(peerID, NewRemoteRecord(peerID, pongPayload.Origin, pongPayload.Latest,
				pongPayload.CheckpointHash, pongPayload.CheckpointActualDiff))
		} else {
			mod.store[peerID].update(pongPayload.Origin, pongPayload.Latest,
				pongPayload.CheckpointHash, pongPayload.CheckpointActualDiff)
		}

	case message.MessageType_REJECT:
		return fmt.Errorf("ping is rejected by remote: %s", string(reply.Payload))
	default:
		return fmt.Errorf("remote replies ping with invalid messgae type: %s", reply.Header.MessageType)
	}

	return nil
}

// getRemoteChainFromLocalLatest just get the remote status from remote
func (mod *syncModule) getRemoteChainFromLocalLatest(record *RemoteRecord) (chain []*ngtypes.Block, err error) {
	latestHash := mod.pow.Chain.GetLatestBlockHash()

	id, s, err := mod.localNode.SendGetChain(record.id, [][]byte{latestHash}, nil) // nil means get MaxBlocks number blocks
	if s == nil {
		return nil, fmt.Errorf("failed to send getchain: %s", err)
	}

	reply, err := wired.ReceiveReply(id, s)
	if err != nil {
		return nil, err
	}

	switch reply.Header.MessageType {
	case message.MessageType_CHAIN:
		chainPayload, err := wired.DecodeChainPayload(reply.Payload)
		if err != nil {
			return nil, fmt.Errorf("failed to send ping: %s", err)
		}

		// TODO: add support for hashes etc
		return chainPayload.Blocks, err

	case message.MessageType_REJECT:
		return nil, fmt.Errorf("getchain is rejected by remote: %s", string(reply.Payload))

	case message.MessageType_NOTFOUND:
		return nil, fmt.Errorf("chain is not found in remote")

	default:
		return nil, fmt.Errorf("remote replies ping with invalid messgae type: %s", reply.Header.MessageType)
	}
}

// getRemoteChain get the chain from remote node
func (mod *syncModule) getRemoteChain(peerID core.PeerID, from [][]byte, to []byte) (chain []*ngtypes.Block, err error) {
	id, s, err := mod.localNode.SendGetChain(peerID, from, to)
	if s == nil {
		return nil, fmt.Errorf("failed to send getchain: %s", err)
	}

	reply, err := wired.ReceiveReply(id, s)
	if err != nil {
		return nil, err
	}

	switch reply.Header.MessageType {
	case message.MessageType_CHAIN:
		chainPayload, err := wired.DecodeChainPayload(reply.Payload)
		if err != nil {
			return nil, fmt.Errorf("failed to send ping: %s", err)
		}

		// TODO: add support for hashes etc
		return chainPayload.Blocks, err

	case message.MessageType_REJECT:
		return nil, fmt.Errorf("getchain is rejected by remote: %s", string(reply.Payload))

	case message.MessageType_NOTFOUND:
		return nil, fmt.Errorf("chain is not found in remote")

	default:
		return nil, fmt.Errorf("remote replies ping with invalid messgae type: %s", reply.Header.MessageType)
	}
}

func (mod *syncModule) getRemoteStateSheet(record *RemoteRecord) (sheet *ngtypes.Sheet, err error) {
	id, s, err := mod.localNode.SendGetSheet(record.id, record.checkpointHeight, record.checkpointHash)
	if s == nil {
		return nil, fmt.Errorf("failed to send getsheet: %s", err)
	}

	reply, err := wired.ReceiveReply(id, s)
	if err != nil {
		return nil, err
	}

	switch reply.Header.MessageType {
	case message.MessageType_SHEET:
		sheetPayload, err := wired.DecodeSheetPayload(reply.Payload)
		if err != nil {
			return nil, fmt.Errorf("failed to send ping: %s", err)
		}

		// TODO: add support for hashes etc
		return sheetPayload.Sheet, err

	case message.MessageType_REJECT:
		return nil, fmt.Errorf("getsheet is rejected by remote: %s", string(reply.Payload))

	case message.MessageType_NOTFOUND:
		return nil, fmt.Errorf("checkpoint's sheet is not found in remote")

	default:
		return nil, fmt.Errorf("remote replies with invalid messgae type: %s", reply.Header.MessageType)
	}
}
