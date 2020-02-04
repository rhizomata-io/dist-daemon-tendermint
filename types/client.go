package types

import abci "github.com/tendermint/tendermint/abci/types"

type Client interface {
	BroadcastTxSync(msg *TxMsg) error
	BroadcastTxAsync(msg *TxMsg) error
	Commit() error
	Query(msg *ViewMsg) (abci.ResponseQuery, error)
}
