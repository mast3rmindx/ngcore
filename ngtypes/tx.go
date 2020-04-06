package ngtypes

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/gogo/protobuf/proto"
	"golang.org/x/crypto/sha3"

	"github.com/cbergoon/merkletree"
	"github.com/mr-tron/base58"

	"github.com/ngchain/ngcore/utils"
)

var (
	ErrTxIsNotSigned         = errors.New("the transaction is not signed")
	ErrTxBalanceInsufficient = errors.New("balance is insufficient for payment")
	ErrTxWrongSign           = errors.New("the signer of transaction is not the own of the account")
)

// NewUnsignedTx will return an unsigned tx, must using Signature()
func NewUnsignedTx(txType TxType, convener uint64, participants [][]byte, values []*big.Int, fee *big.Int, nonce uint64, extraData []byte) *Tx {
	header := &TxHeader{
		Version:      Version,
		Type:         txType,
		Convener:     convener,
		Participants: participants,
		Fee:          fee.Bytes(),
		Values:       BigIntsToBytesList(values),
		Nonce:        nonce,
		Extra:        extraData,
	}

	return &Tx{
		Header: header,

		R: nil,
		S: nil,
	}
}

// IsSigned will return whether the op has been signed
func (m *Tx) IsSigned() bool {
	if m.R == nil || m.S == nil {
		return false
	}
	return true
}

// Verify helps verify the transaction whether signed by the public key owner
func (m *Tx) Verify(pubKey ecdsa.PublicKey) error {
	if m.R == nil || m.S == nil {
		log.Panic("unsigned transaction")
	}

	b, err := proto.Marshal(m.Header)
	if err != nil {
		return err
	}

	hash := sha3.Sum256(b)
	if !ecdsa.Verify(&pubKey, hash[:], new(big.Int).SetBytes(m.R), new(big.Int).SetBytes(m.S)) {
		return ErrTxWrongSign
	}

	return nil
}

// Bs58 is a tx's ReadableID in string
func (m *Tx) Bs58() string {
	b, err := proto.Marshal(m)
	if err != nil {
		log.Error(err)
	}
	return base58.FastBase58Encoding(b)
}

// HashHex is a tx's ReadableID in string
func (m *Tx) HashHex() string {
	b, err := m.CalculateHash()
	if err != nil {
		log.Error(err)
		return ""
	}

	return hex.EncodeToString(b)
}

// CalculateHash mainly for calculating the tire root of txs and sign tx
func (m *Tx) CalculateHash() ([]byte, error) {
	raw, err := m.Marshal()
	if err != nil {
		log.Error(err)
	}

	hash := sha3.Sum256(raw)
	return hash[:], nil
}

// Equals mainly for calculating the tire root of txs
func (m *Tx) Equals(other merkletree.Content) (bool, error) {
	tx, ok := other.(*Tx)
	if !ok {
		return false, errors.New("invalid transaction type")
	}

	otherHash, err := tx.Header.CalculateHash()
	if err != nil {
		return false, err
	}
	mHash, err := m.Header.CalculateHash()
	if err != nil {
		return false, err
	}

	return bytes.Equal(otherHash, mHash), nil
}

// TxsToMerkleTreeContents make a []merkletree.Content whose values is from txs
func TxsToMerkleTreeContents(txs []*Tx) []merkletree.Content {
	mtc := make([]merkletree.Content, len(txs))
	for i := range txs {
		mtc[i] = txs[i]
	}

	return mtc
}

func (m *Tx) Copy() *Tx {
	tx := proto.Clone(m).(*Tx)
	return tx
}

// BigIntsToBytesList is a helper converts bigInts to raw bytes slice
func BigIntsToBytesList(bigInts []*big.Int) [][]byte {
	bytesList := make([][]byte, len(bigInts))
	for i := 0; i < len(bigInts); i++ {
		bytesList[i] = bigInts[i].Bytes()
	}
	return bytesList
}

func (m *Tx) CheckGenerate() error {
	if m.Header == nil {
		return errors.New("generate is missing header")
	}

	if m.GetConvener() != 0 {
		return fmt.Errorf("generate's convener should be 0")
	}

	if len(m.GetValues()) != len(m.GetParticipants()) {
		return fmt.Errorf("generate should have same len with participants")
	}

	if !bytes.Equal(m.TotalCharge().Bytes(), OneBlockReward.Bytes()) {
		return fmt.Errorf("wrong block reward")
	}

	if !bytes.Equal(m.GetFee(), GetBig0Bytes()) {
		return fmt.Errorf("generate's fee should be ZERO")
	}

	publicKey := utils.Bytes2ECDSAPublicKey(m.GetParticipants()[0])
	if err := m.Verify(publicKey); err != nil {
		return err
	}

	return nil
}

func (m *Tx) CheckRegister() error {
	if m.Header == nil {
		return errors.New("register is missing header")
	}

	if m.GetConvener() != 0 {
		return fmt.Errorf("register's convener should be 0")
	}

	if len(m.GetParticipants()) != 1 {
		return fmt.Errorf("register should have only one participant")
	}

	if len(m.GetValues()) != 1 {
		return fmt.Errorf("register should have only one value")
	}

	if !bytes.Equal(m.GetValues()[0], GetBig0Bytes()) {
		return fmt.Errorf("register should have only one 0 value")
	}

	if new(big.Int).SetBytes(m.GetFee()).Cmp(new(big.Int).Mul(NG, big.NewInt(10))) < 0 {
		return fmt.Errorf("register should have at least 10NG fee")
	}

	if len(m.GetExtra()) != 8 {
		return fmt.Errorf("register should have uint64 little-endian bytes as extra")
	}

	publicKey := utils.Bytes2ECDSAPublicKey(m.GetParticipants()[0])
	if err := m.Verify(publicKey); err != nil {
		return err
	}

	return nil
}

func (m *Tx) CheckLogout(key ecdsa.PublicKey) error {
	if m.Header == nil {
		return errors.New("logout is missing header")
	}

	if len(m.GetParticipants()) != 0 {
		return fmt.Errorf("logout should have NO participant")
	}

	if m.GetConvener() == 0 {
		return fmt.Errorf("logout's convener should NOT be 0")
	}

	if len(m.GetValues()) != 0 {
		return fmt.Errorf("logout should have NO value")
	}

	if len(m.GetValues()) != len(m.GetParticipants()) {
		return fmt.Errorf("logout should have same len with participants")
	}

	if err := m.Verify(key); err != nil {
		return err
	}

	return nil
}

func (m *Tx) CheckTransaction(key ecdsa.PublicKey) error {
	if m.Header == nil {
		return errors.New("transaction is missing header")
	}

	if m.GetConvener() == 0 {
		return fmt.Errorf("transaction's convener should NOT be 0")
	}

	if len(m.GetValues()) != len(m.GetParticipants()) {
		return fmt.Errorf("transaction should have same len with participants")
	}

	if err := m.Verify(key); err != nil {
		return err
	}

	return nil
}

func (m *Tx) CheckAssign(key ecdsa.PublicKey) error {
	if m.Header == nil {
		return errors.New("assign is missing header")
	}

	if m.GetConvener() == 0 {
		return fmt.Errorf("assign's convener should NOT be 0")
	}

	if len(m.GetParticipants()) != 0 {
		return fmt.Errorf("assign should have NO participant")
	}

	if len(m.GetValues()) != 0 {
		return fmt.Errorf("assign should have NO value")
	}

	if err := m.Verify(key); err != nil {
		return err
	}

	return nil
}

func (m *Tx) CheckAppend(key ecdsa.PublicKey) error {
	if m.Header == nil {
		return errors.New("append is missing header")
	}

	if len(m.GetParticipants()) != 0 {
		return fmt.Errorf("append should have NO participant")
	}

	if m.GetConvener() == 0 {
		return fmt.Errorf("append's convener should NOT be 0")
	}

	if len(m.GetValues()) != 0 {
		return fmt.Errorf("append should have NO value")
	}

	if err := m.Verify(key); err != nil {
		return err
	}

	return nil
}

// Signature will re-sign the Tx with private key
func (m *Tx) Signature(privKey *ecdsa.PrivateKey) (err error) {
	b, err := proto.Marshal(m.Header)
	if err != nil {
		log.Error(err)
	}

	hash := sha3.Sum256(b)
	r, s, err := ecdsa.Sign(rand.Reader, privKey, hash[:])
	if err != nil {
		log.Panic(err)
	}

	m.R = r.Bytes()
	m.S = s.Bytes()

	return
}

func (m *Tx) GetType() TxType {
	return m.Header.GetType()
}

func (m *Tx) GetConvener() uint64 {
	return m.Header.GetConvener()
}

func (m *Tx) GetValues() [][]byte {
	return m.Header.GetValues()
}

func (m *Tx) GetParticipants() [][]byte {
	return m.Header.GetParticipants()
}

func (m *Tx) GetFee() []byte {
	return m.Header.GetFee()
}

func (m *Tx) GetNonce() uint64 {
	return m.Header.GetNonce()
}

func (m *Tx) GetVersion() int32 {
	return m.Header.GetVersion()
}

func (m *Tx) GetExtra() []byte {
	return m.Header.GetExtra()
}

func (m *Tx) TotalCharge() *big.Int {
	return m.Header.TotalCharge()
}

// GetGenesisGeneration is a constructed function
func GetGenesisGeneration() *Tx {
	gen := NewUnsignedTx(
		TX_GENERATION,
		0,
		[][]byte{GenesisPK},
		[]*big.Int{OneBlockReward},
		GetBig0(),
		1,
		nil,
	)

	// FIXME: before init network should manually init the R & S
	gen.R, _ = hex.DecodeString("e60f90c8e8bb717cf30cf59a8bec8d17f189a5a4e0a621f4c2ce2d24a0443d1f")
	gen.S, _ = hex.DecodeString("dfcba58223e100569991a856ca139287714e5cd53074bc962328a602fe3b81bf")

	return gen
}
