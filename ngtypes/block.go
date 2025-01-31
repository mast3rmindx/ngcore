package ngtypes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"runtime"
	"sync"
	"time"

	logging "github.com/ipfs/go-log/v2"
	"github.com/ngchain/go-randomx"

	"google.golang.org/protobuf/proto"

	"golang.org/x/crypto/sha3"

	"github.com/ngchain/ngcore/utils"
)

var log = logging.Logger("types")

// IsUnsealing checks whether the block is unsealing.
func (x *Block) IsUnsealing() bool {
	return x.GetTrieHash() != nil
}

// IsSealed checks whether the block is sealed.
func (x *Block) IsSealed() bool {
	return x.GetNonce() != nil
}

// IsHead will check whether the Block is the head(checkpoint).
func (x *Block) IsHead() bool {
	return x.GetHeight()%BlockCheckRound == 0
}

// IsTail will check whether the Block is the tail(the one before head).
func (x *Block) IsTail() bool {
	return (x.GetHeight()+1)%BlockCheckRound == 0
}

// IsGenesis will check whether the Block is the genesis block.
func (x *Block) IsGenesis() bool {
	return bytes.Equal(x.Hash(), GetGenesisBlockHash(x.Network))
}

// GetPoWRawHeader will return a complete raw for block hash.
// When nonce is not nil, the RawHeader will use the nonce param not the x.Nonce.
func (x *Block) GetPoWRawHeader(nonce []byte) []byte {
	//lenRaw := 1 + // network size
	//	HeightSize+
	//	TimestampSize +
	//	HashSize +
	//	HashSize + // unknown length
	//	HashSize +
	//	NonceSize
	raw := make([]byte, 121)

	raw[0] = byte(x.Network)
	binary.LittleEndian.PutUint64(raw[1:], x.Height)
	binary.LittleEndian.PutUint64(raw[9:17], uint64(x.GetTimestamp()))
	copy(raw[17:49], x.GetPrevBlockHash())
	copy(raw[49:81], x.GetTrieHash())
	copy(raw[81:113], utils.ReverseBytes(x.GetDifficulty())) // uint256

	if nonce == nil {
		copy(raw[113:121], x.GetNonce())
	} else {
		copy(raw[113:121], nonce)
	}

	return raw
}

// ApplyPoWRawAndTxs will apply the raw pow of header and txs to the block.
func (x *Block) ApplyPoWRawAndTxs(raw []byte, txs []*Tx) error {
	//lenRaw := 1 + // network size
	//	HeightSize+
	//	TimestampSize +
	//	HashSize +
	//	HashSize + // unknown length
	//	HashSize +
	//	NonceSize
	if len(raw) != 121 {
		return fmt.Errorf("wrong length of PoW raw bytes")
	}

	*x = Block{
		Network:       NetworkType(raw[0]),
		Height:        binary.LittleEndian.Uint64(raw[1:9][:]),
		Timestamp:     int64(binary.LittleEndian.Uint64(raw[9:17])),
		PrevBlockHash: raw[17:49],
		TrieHash:      raw[49:81],
		Difficulty:    bytes.TrimLeft(utils.ReverseBytes(raw[81:113]), string(byte(0))), // remove left padding
		Nonce:         raw[113:121],
		Txs:           txs,
	}

	return nil
}

// PowHash will help you get the pow hash of block.
func (x *Block) PowHash() []byte {
	cache, err := randomx.AllocCache(randomx.FlagJIT)
	if err != nil {
		panic(err)
	}
	defer randomx.ReleaseCache(cache)

	randomx.InitCache(cache, x.PrevBlockHash)
	ds, err := randomx.AllocDataset(randomx.FlagJIT)
	if err != nil {
		panic(err)
	}
	defer randomx.ReleaseDataset(ds)

	count := randomx.DatasetItemCount()
	var wg sync.WaitGroup
	var workerNum = uint32(runtime.NumCPU())
	for i := uint32(0); i < workerNum; i++ {
		wg.Add(1)
		a := (count * i) / workerNum
		b := (count * (i + 1)) / workerNum
		go func() {
			defer wg.Done()
			randomx.InitDataset(ds, cache, a, b-a)
		}()
	}
	wg.Wait()

	vm, err := randomx.CreateVM(cache, ds, randomx.FlagJIT)
	if err != nil {
		panic(err)
	}
	defer randomx.DestroyVM(vm)

	return randomx.CalculateHash(vm, x.GetPoWRawHeader(nil))
}

// ToUnsealing converts a bare block to an unsealing block.
func (x *Block) ToUnsealing(txsWithGen []*Tx) (*Block, error) {
	b := proto.Clone(x).(*Block)

	if txsWithGen[0].GetType() != TxType_GENERATE {
		return nil, fmt.Errorf("first tx shall be a generate")
	}

	for i := 1; i < len(txsWithGen); i++ {
		if txsWithGen[i].GetType() == TxType_GENERATE {
			return nil, fmt.Errorf("except first, other tx shall not be a generate")
		}
	}

	b.TrieHash = NewTxTrie(txsWithGen).TrieRoot()
	b.Txs = txsWithGen

	return b, nil
}

// ToSealed converts an unsealing block to a sealed block.
func (x *Block) ToSealed(nonce []byte) (*Block, error) {
	if !x.IsUnsealing() {
		return nil, fmt.Errorf("the block is bare")
	}

	if len(nonce) != NonceSize {
		return nil, fmt.Errorf("nonce length %d is incorrect", len(nonce))
	}

	x.Nonce = nonce

	return x, nil
}

// verifyNonce will verify whether the nonce meets the target.
func (x *Block) verifyNonce() error {
	diff := new(big.Int).SetBytes(x.GetDifficulty())
	target := new(big.Int).Div(MaxTarget, diff)

	if new(big.Int).SetBytes(x.PowHash()).Cmp(target) < 0 {
		return nil
	}

	return fmt.Errorf("block@%d's nonce %x is invalid", x.GetHeight(), x.GetNonce())
}

// GetActualDiff returns the diff decided by nonce.
func (x *Block) GetActualDiff() *big.Int {
	return new(big.Int).Div(MaxTarget, new(big.Int).SetBytes(x.PowHash()))
}

// NewBareBlock will return an unsealing block and
// then you need to add txs and seal with the correct N.
func NewBareBlock(network NetworkType, height uint64, blockTime int64, prevBlockHash []byte, diff *big.Int) *Block {
	return &Block{
		Network:       network,
		Height:        height,
		Timestamp:     blockTime,
		PrevBlockHash: prevBlockHash,
		TrieHash:      make([]byte, HashSize),

		Difficulty: diff.Bytes(),
		Nonce:      make([]byte, NonceSize),

		Txs: make([]*Tx, 0),
	}
}

// CheckError will check the errors in block inner fields.
func (x *Block) CheckError() error {
	//if x.Network != Network {
	//	return fmt.Errorf("block's network id is incorrect")
	//}
	// DONE: do network check on consensus

	if len(x.PrevBlockHash) != HashSize {
		return fmt.Errorf("block%d's PrevBlockHash length is incorrect", x.GetHeight())
	}

	if len(x.TrieHash) != HashSize {
		return fmt.Errorf("block%d's TrieHash length is incorrect", x.GetHeight())
	}

	if len(x.Nonce) != NonceSize {
		return fmt.Errorf("block%d's Nonce length is incorrect", x.GetHeight())
	}

	if x.Timestamp > time.Now().Unix() {
		return fmt.Errorf("block%d's timestamp %d is invalid", x.GetHeight(), x.Timestamp)
	}

	if !x.IsSealed() {
		return fmt.Errorf("block@%d has not sealed with nonce", x.GetHeight())
	}

	if !bytes.Equal(NewTxTrie(x.Txs).TrieRoot(), x.GetTrieHash()) {
		return fmt.Errorf("the merkle tree in block@%d is invalid", x.GetHeight())
	}

	err := x.verifyNonce()
	if err != nil {
		return err
	}

	err = x.verifyNonce()
	if err != nil {
		return err
	}

	return nil
}

func (x *Block) verifyHash() error {
	if x.Id == nil {
		x.Hash()
		return nil
	}

	b := proto.Clone(x).(*Block)
	b.Txs = nil // txs can be represented by triehash

	raw, err := utils.Proto.Marshal(b)
	if err != nil {
		panic(err)
	}

	hash := sha3.Sum256(raw)

	if !bytes.Equal(hash[:], x.Id) {
		return fmt.Errorf("block@%d hash %x not match its id %x", x.Height, hash, x.Id)
	}

	return nil
}

// Hash will help you get the hash of block.
func (x *Block) Hash() []byte {
	if x.Id == nil {
		b := proto.Clone(x).(*Block)
		b.Txs = nil // txs can be represented by triehash

		raw, err := utils.Proto.Marshal(b)
		if err != nil {
			panic(err)
		}

		hash := sha3.Sum256(raw)

		x.Id = hash[:]
	}

	return x.Id
}

// GetPrevHash is a helper to get the prev block hash from block header.
func (x *Block) GetPrevHash() []byte {
	return x.GetPrevBlockHash()
}

// GetGenesisBlock will return a complete sealed GenesisBlock.
func GetGenesisBlock(network NetworkType) *Block {
	txs := []*Tx{
		GetGenesisGenerateTx(network),
	}

	genesisBlock := &Block{
		Network:   network,
		Height:    0,
		Timestamp: GetGenesisTimestamp(network),

		PrevBlockHash: make([]byte, HashSize),
		TrieHash:      NewTxTrie(txs).TrieRoot(),

		Difficulty: minimumBigDifficulty.Bytes(), // this is a number, dont put any padding on
		Nonce:      GetGenesisBlockNonce(network),
		Txs:        txs,
	}

	return genesisBlock
}

func GetGenesisBlockHash(network NetworkType) []byte {
	return GetGenesisBlock(network).Hash()
}
