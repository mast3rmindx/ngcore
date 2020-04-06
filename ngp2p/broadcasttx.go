package ngp2p

import (
	"context"

	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/ngchain/ngcore/ngtypes"
)

func (b *Broadcaster) broadcastTx(tx *ngtypes.Tx) bool {
	raw, err := tx.Marshal()
	if err != nil {
		log.Errorf("failed to sign pb data")
		return false
	}

	err = b.topics[broadcastBlockTopic].Publish(context.Background(), raw)
	if err != nil {
		log.Error(err)
		return false
	}

	return true
}

func (b *Broadcaster) onBroadcastTx(msg *pubsub.Message) {
	var tx = &ngtypes.Tx{}
	err := tx.Unmarshal(msg.Data)
	if err != nil {
		log.Error(err)
		return
	}

	err = b.node.txPool.PutTxs(tx)
	if err != nil {
		log.Error(err)
		return
	}
}
