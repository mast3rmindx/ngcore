package consensus

import (
	"crypto/ecdsa"
	"github.com/ngin-network/ngcore/chain"
	miner2 "github.com/ngin-network/ngcore/consensus/miner"
	"github.com/ngin-network/ngcore/ngtypes"
	"github.com/ngin-network/ngcore/sheetManager"
	"github.com/ngin-network/ngcore/txpool"
	"sync"
)

// the pow
type Consensus struct {
	sync.RWMutex

	template     *ngtypes.Block
	SheetManager *sheetManager.SheetManager

	privateKey *ecdsa.PrivateKey
	Chain      *chain.Chain

	TxPool *txpool.TxPool

	mining bool
	miner  *miner2.Miner
}

func NewConsensusManager(mining bool) *Consensus {
	return &Consensus{
		template:     nil,
		SheetManager: nil,
		privateKey:   nil,
		Chain:        nil,
		TxPool:       nil,

		mining: mining,
	}
}

func (c *Consensus) Init(chain *chain.Chain, sheetManager *sheetManager.SheetManager, privateKey *ecdsa.PrivateKey, txPool *txpool.TxPool) {
	c.privateKey = privateKey
	c.SheetManager = sheetManager
	c.Chain = chain
	c.TxPool = txPool
}
