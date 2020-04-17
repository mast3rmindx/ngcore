package txpool

import (
	"testing"

	"github.com/ngchain/ngcore/ngsheet"
	"github.com/ngchain/ngcore/ngtypes"
)

func TestNewTxPool(t *testing.T) {
	sheetManager := ngsheet.NewSheetManager()
	pool := GetTxPool(sheetManager)

	sheetManager.Init(ngtypes.GetGenesisBlock())
	pool.Init(nil)
}
