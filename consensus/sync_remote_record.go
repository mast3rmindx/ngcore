package consensus

import (
	"bytes"
	"math/big"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/ngchain/ngcore/ngtypes"
	"github.com/ngchain/ngcore/storage"
)

type remoteRecord struct {
	id                   peer.ID
	origin               uint64 // rank
	latest               uint64
	checkpointHash       []byte   // trigger
	checkpointActualDiff *big.Int // rank
	lastChatTime         int64
}

// RULE: checkpoint fork: when a node mined a checkpoint, all other node are forced to start sync
func (r *remoteRecord) shouldSync() bool {
	if r.latest > storage.GetChain().GetLatestBlockHeight() {
		return true
	}

	return false
}

// RULE: when forking?
// Situation #1: remote height is higher than local, AND checkpoint is on higher level
// Situation #2: remote height is higher than local, AND checkpoint is on same level, AND remote checkpoint takes more rank (with more ActualDiff)
func (r *remoteRecord) shouldFork() bool {
	cp := storage.GetChain().GetLatestCheckpoint()
	cpHash, _ := cp.CalculateHash()

	h := storage.GetChain().GetLatestBlockHeight()

	if !bytes.Equal(r.checkpointHash, cpHash) &&
		r.latest%ngtypes.BlockCheckRound > h%ngtypes.BlockCheckRound {
		return true
	}

	if !bytes.Equal(r.checkpointHash, cpHash) &&
		r.latest > h &&
		r.latest%ngtypes.BlockCheckRound == h%ngtypes.BlockCheckRound &&
		r.checkpointActualDiff.Cmp(cp.GetActualDiff()) > 0 {
		return true
	}

	return false
}