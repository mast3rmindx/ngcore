package txpool

import (
	"github.com/ngchain/ngcore/ngtypes"
)

func (p *TxPool) DelBlockTxs(txs ...*ngtypes.Transaction) {
	p.Lock()
	defer p.Unlock()

	for i := range txs {
		if p.Queuing[txs[i].GetConvener()] != nil {
			delete(p.Queuing[txs[i].GetConvener()], txs[i].GetNonce())
		}
	}
}

// CheckTxs will check txs self and error in sheet
func (p *TxPool) CheckTxs(txs ...*ngtypes.Transaction) bool {
	if p.sheetManager.CheckTxs(txs...) != nil {
		return false
	}

	for i := range txs {
		if txs[i].Check() != nil {
			return false
		}
	}

	return true
}