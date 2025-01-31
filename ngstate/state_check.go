package ngstate

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/ngchain/ngcore/utils"

	"github.com/dgraph-io/badger/v2"

	"github.com/ngchain/ngcore/ngtypes"
)

// CheckBlockTxs will check all requirements for txs in block
func CheckBlockTxs(txn *badger.Txn, block *ngtypes.Block) error {
	for i := 0; i < len(block.Txs); i++ {
		tx := block.Txs[i]
		// check tx is signed
		if !tx.IsSigned() {
			return fmt.Errorf("tx is not signed")
		}

		// check the tx's extra size is necessary
		if len(tx.Extra) > ngtypes.TxMaxExtraSize {
			return fmt.Errorf("tx is too large")
		}

		switch tx.GetType() {
		case ngtypes.TxType_GENERATE: // generate
			if err := checkGenerate(txn, tx, block.Height); err != nil {
				return err
			}

		case ngtypes.TxType_REGISTER: // register
			if err := checkRegister(txn, tx); err != nil {
				return err
			}

		case ngtypes.TxType_LOGOUT: // logout
			if err := checkLogout(txn, tx); err != nil {
				return err
			}

		case ngtypes.TxType_TRANSACT: // transaction
			if err := checkTransaction(txn, tx); err != nil {
				return err
			}

		case ngtypes.TxType_APPEND: // append
			if err := checkAppend(txn, tx); err != nil {
				return err
			}

		case ngtypes.TxType_DELETE: // delete
			if err := checkDelete(txn, tx); err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid tx type")
		}
	}

	return nil
}

// CheckTx will check the requirements for one tx (except generate tx)
func CheckTx(txn *badger.Txn, tx *ngtypes.Tx) error {
	// check tx is signed
	if !tx.IsSigned() {
		return fmt.Errorf("tx is not signed")
	}

	// check the tx's extra size is necessary
	if len(tx.Extra) > ngtypes.TxMaxExtraSize {
		return fmt.Errorf("tx is too large")
	}

	switch tx.GetType() {
	case ngtypes.TxType_GENERATE: // generate
		return fmt.Errorf("cannot check generate tx with CheckTx")

	case ngtypes.TxType_REGISTER: // register
		if err := checkRegister(txn, tx); err != nil {
			return err
		}

	case ngtypes.TxType_LOGOUT: // logout
		if err := checkLogout(txn, tx); err != nil {
			return err
		}

	case ngtypes.TxType_TRANSACT: // transact
		if err := checkTransaction(txn, tx); err != nil {
			return err
		}

	case ngtypes.TxType_DELETE: // delete
		if err := checkDelete(txn, tx); err != nil {
			return err
		}

	case ngtypes.TxType_APPEND: // append
		if err := checkAppend(txn, tx); err != nil {
			return err
		}
	}

	return nil
}

// checkGenerate checks the generate tx
func checkGenerate(txn *badger.Txn, generateTx *ngtypes.Tx, blockHeight uint64) error {

	item, err := txn.Get(append(numToAccountPrefix, ngtypes.AccountNum(generateTx.GetConvener()).Bytes()...))
	if err != nil {
		return fmt.Errorf("cannot find convener %d: %s", generateTx.GetConvener(), err)
	}

	rawConvener, err := item.ValueCopy(nil)
	if err != nil {
		return fmt.Errorf("cannot get convener account %d: %s", generateTx.GetConvener(), err)
	}

	convener := new(ngtypes.Account)
	err = utils.Proto.Unmarshal(rawConvener, convener)
	if err != nil {
		return err
	}

	// check structure and key
	if err := generateTx.CheckGenerate(blockHeight); err != nil {
		return err
	}

	// check rew

	return nil
}

// checkRegister checks the register tx
func checkRegister(txn *badger.Txn, registerTx *ngtypes.Tx) error {
	// check structure and key
	if err := registerTx.CheckRegister(); err != nil {
		return err
	}

	// check balance
	payerAddr := registerTx.GetParticipants()[0]
	payerBalance, err := getBalance(txn, payerAddr)
	if err != nil {
		return err
	}

	expenditure := registerTx.TotalExpenditure()
	if payerBalance.Cmp(expenditure) < 0 {
		return fmt.Errorf("balance is insufficient for register")
	}

	// check existing ownership
	if addrHasAccount(txn, payerAddr) {
		return fmt.Errorf("one address cannot repeat registering accounts")
	}

	// check newAccountNum
	newAccountNum := binary.LittleEndian.Uint64(registerTx.GetExtra())
	if accountNumExists(txn, ngtypes.AccountNum(newAccountNum)) {
		return fmt.Errorf("failed to register account@%d, account is already used by others", newAccountNum)
	}

	return nil
}

// checkLogout checks logout tx
func checkLogout(txn *badger.Txn, logoutTx *ngtypes.Tx) error {
	convener, err := getAccountByNum(txn, ngtypes.AccountNum(logoutTx.GetConvener()))
	if err != nil {
		return err
	}

	// check structure and key
	if err = logoutTx.CheckLogout(ngtypes.Address(convener.Owner).PubKey()); err != nil {
		return err
	}

	// check balance
	totalCharge := logoutTx.TotalExpenditure()
	convenerBalance, err := getBalance(txn, convener.Owner)
	if err != nil {
		return err
	}

	if convenerBalance.Cmp(totalCharge) < 0 {
		return fmt.Errorf("balance is insufficient for logout")
	}

	if len(convener.Contract) != 0 {
		return fmt.Errorf("you should clear your contract before logout")
	}

	if len(convener.Context) != 0 {
		return fmt.Errorf("you should clear your context before logout")
	}

	return nil
}

// checkTransaction checks normal transaction tx
func checkTransaction(txn *badger.Txn, transactionTx *ngtypes.Tx) error {
	convener, err := getAccountByNum(txn, ngtypes.AccountNum(transactionTx.Convener))
	if err != nil {
		return err
	}

	// check structure and key
	if err = transactionTx.CheckTransaction(ngtypes.Address(convener.Owner).PubKey()); err != nil {
		return err
	}

	// check balance
	totalCharge := transactionTx.TotalExpenditure()
	convenerBalance, err := getBalance(txn, convener.Owner)
	if err != nil {
		return err
	}

	if convenerBalance.Cmp(totalCharge) < 0 {
		return fmt.Errorf("balance is insufficient for transaction")
	}

	return nil
}

// checkAppend checks append tx
func checkAppend(txn *badger.Txn, appendTx *ngtypes.Tx) error {
	convener, err := getAccountByNum(txn, ngtypes.AccountNum(appendTx.Convener))
	if err != nil {
		return err
	}

	// check structure and key
	if err = appendTx.CheckAppend(ngtypes.Address(convener.Owner).PubKey()); err != nil {
		return err
	}

	// check balance
	totalCharge := appendTx.TotalExpenditure()
	convenerBalance, err := getBalance(txn, convener.Owner)
	if err != nil {
		return err
	}

	if convenerBalance.Cmp(totalCharge) < 0 {
		return fmt.Errorf("balance is insufficient for append")
	}

	var appendExtra ngtypes.AppendExtra
	err = proto.Unmarshal(appendTx.Extra, &appendExtra)
	if err != nil {
		return err
	}

	if appendExtra.Pos >= uint64(len(convener.Contract)) {
		return fmt.Errorf("append pos is out of bound")
	}

	return nil
}

// checkDelete checks delete tx
func checkDelete(txn *badger.Txn, deleteTx *ngtypes.Tx) error {
	convener, err := getAccountByNum(txn, ngtypes.AccountNum(deleteTx.Convener))
	if err != nil {
		return err
	}

	// check structure and key
	if err = deleteTx.CheckDelete(ngtypes.Address(convener.Owner).PubKey()); err != nil {
		return err
	}

	// check balance
	totalCharge := deleteTx.TotalExpenditure()
	convenerBalance, err := getBalance(txn, convener.Owner)
	if err != nil {
		return err
	}

	if convenerBalance.Cmp(totalCharge) < 0 {
		return fmt.Errorf("balance is insufficient for delete")
	}

	var appendExtra ngtypes.DeleteExtra
	err = proto.Unmarshal(deleteTx.Extra, &appendExtra)
	if err != nil {
		return err
	}

	if appendExtra.Pos >= uint64(len(convener.Contract)) {
		return fmt.Errorf("delete pos is out of bound")
	}

	if appendExtra.Pos+uint64(len(appendExtra.Content)) >= uint64(len(convener.Contract)) {
		return fmt.Errorf("delete content length is out of bound")
	}

	if !bytes.Equal(
		convener.Contract[int(appendExtra.Pos):int(appendExtra.Pos)+len(appendExtra.Content)],
		appendExtra.Content) {
		return fmt.Errorf("delete content length is invalid")
	}

	return nil
}
