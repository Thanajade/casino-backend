package main

import (
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/rs/zerolog/log"
)

func NewSigndice(contract, casinoAccount string, requestID uint64, signature string) *eos.Action {
	return &eos.Action{
		Account: eos.AN(contract),
		Name:    eos.ActN("sgdicesecond"),
		Authorization: []eos.PermissionLevel{
			{Actor: eos.AN(casinoAccount), Permission: eos.PN("signidice")},
		},
		ActionData: eos.NewActionData(Signidice{
			requestID,
			signature,
		}),
	}
}

// Game contract's sgdicesecond action parameters
type Signidice struct {
	RequestID uint64 `json:"req_id"`
	Signature string `json:"sign"`
}


func GetSigndiceTransaction(api *eos.API, contract, casinoAccount string, requestID uint64, signature ecc.Signature) (*eos.SignedTransaction, *eos.PackedTransaction) {
	action := NewSigndice(contract, casinoAccount, requestID, string(signature.Content))
	txOpts := &eos.TxOptions{}

	if err := txOpts.FillFromChain(api); err != nil {
		log.Error().Msgf("filling tx opts: %s", err.Error())
		return nil, nil
	}
	tx := eos.NewTransaction([]*eos.Action{action}, txOpts)
	signedTx, packedTx, err := api.SignTransaction(tx, txOpts.ChainID, eos.CompressionNone)
	if err != nil {
		log.Error().Msgf("sign transaction: %s", err.Error())
		return nil, nil
	}
	return signedTx, packedTx
}
