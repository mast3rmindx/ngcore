package jsonrpc

import (
	"math/big"

	"github.com/c0mm4nd/go-jsonrpc2"
	"github.com/mr-tron/base58"
	"github.com/ngchain/secp256k1"

	"github.com/ngchain/ngcore/ngtypes"
	"github.com/ngchain/ngcore/utils"
)

// some utils for wallet clients

type getAddressParams struct {
	PrivateKeys []string
}

type getAddressReply struct {
	Address ngtypes.Address
}

// publicKeyToAddressFunc helps client to get the schnorr publickey of private keys
func (s *Server) publicKeyToAddressFunc(msg *jsonrpc2.JsonRpcMessage) *jsonrpc2.JsonRpcMessage {
	var params getAddressParams
	err := utils.JSON.Unmarshal(*msg.Params, &params)
	if err != nil {
		log.Error(err)
		return jsonrpc2.NewJsonRpcError(msg.ID, jsonrpc2.NewError(0, err))
	}

	privKeys := make([]*secp256k1.PrivateKey, len(params.PrivateKeys))
	for i := 0; i < len(params.PrivateKeys); i++ {
		bPriv, err := base58.FastBase58Decoding(params.PrivateKeys[i])
		if err != nil {
			log.Error(err)
			return jsonrpc2.NewJsonRpcError(msg.ID, jsonrpc2.NewError(0, err))
		}

		privKeys[i] = secp256k1.NewPrivateKey(new(big.Int).SetBytes(bPriv))
	}

	addr, err := ngtypes.NewAddressFromMultiKeys(privKeys...)
	if err != nil {
		log.Error(err)
		return jsonrpc2.NewJsonRpcError(msg.ID, jsonrpc2.NewError(0, err))
	}

	result := getAddressReply{
		Address: addr,
	}

	raw, err := utils.JSON.Marshal(result)
	if err != nil {
		log.Error(err)
		return jsonrpc2.NewJsonRpcError(msg.ID, jsonrpc2.NewError(0, err))
	}

	return jsonrpc2.NewJsonRpcSuccess(msg.ID, raw)
}
