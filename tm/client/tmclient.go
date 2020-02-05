package client

import (
	"errors"
	
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/rpc/core"
	rpctypes "github.com/tendermint/tendermint/rpc/lib/types"
	
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
)

type TMClient struct {
	config *cfg.Config
	logger log.Logger
	codec *types.Codec
}

var _ types.Client = (*TMClient)(nil)


func NewClient( config *cfg.Config, logger log.Logger, codec *types.Codec) types.Client {
	return &TMClient{ config, logger, codec}
}

func (client *TMClient) BroadcastTxSync(msg *types.TxMsg) (err error) {
	msgBytes, err := types.EncodeTxMsg(msg)
	if err != nil {
		client.logger.Error("MarshalBinaryBare : Msg : ", err)
		return err
	}
	
	_, err = core.BroadcastTxSync(&rpctypes.Context{}, msgBytes)
	if err != nil {
		client.logger.Error("TMClient#BroadcastTxSync ", err)
		return err
	}
	return err
}

func (client *TMClient) BroadcastTxAsync(msg *types.TxMsg) (err error) {
	msgBytes, err := types.EncodeTxMsg(msg)
	if err != nil {
		client.logger.Error("TMClient#BroadcastTxAsync : EncodeTxMsg", err)
		return err
	}
	
	_, err = core.BroadcastTxAsync(&rpctypes.Context{}, msgBytes)
	if err != nil {
		client.logger.Error("TMClient#BroadcastTxAsync", err)
		return err
	}
	return err
}

func (client *TMClient) Commit() (err error) {
	_, err = core.Commit(&rpctypes.Context{}, nil)
	if err != nil {
		client.logger.Error("Commit : ", err)
	}
	return err
}


func (client *TMClient) Query(msg *types.ViewMsg) (data []byte, err error) {
	msgBytes, err := types.EncodeViewMsg(msg)
	if err != nil {
		client.logger.Error("TMClient#Query : EncodeViewMsg", err)
		return data, err
	}
	
	res, err := core.ABCIQuery(&rpctypes.Context{}, msg.Path, bytes.HexBytes(msgBytes), 0, false)
	
	if err != nil {
		client.logger.Error("TMClient#Query : ABCIQuery", err)
	}
	
	return res.Response.Value, err
}


func (client *TMClient) GetObject(msg *types.ViewMsg, obj interface{}) (err error){
	data, err := client.Query(msg)
	
	if err != nil {
		client.logger.Error("TMClient#GetObject : ", err)
		return err
	}
	
	err = client.UnmarshalObject(data, obj)
	
	if err != nil {
		client.logger.Error("TMClient#GetObject Unmarshal : ", err)
		return err
	}
	return err
}

func (client *TMClient) GetMany(msg *types.ViewMsg, handler types.KVHandler) (err error){
	data, err := client.Query(msg)
	
	if err != nil {
		client.logger.Error("TMClient#GetMany : ", err)
		return err
	}
	
	kvArray := []types.KeyValue{}
	err = client.UnmarshalObject(data, &kvArray)
	
	if err != nil {
		client.logger.Error("TMClient#GetMany Unmarshal : ", err)
		return err
	}
	
	for _, kv := range kvArray {
		if !handler(kv.Key, kv.Value) {
			break
		}
	}
	return err
}

func (client *TMClient) UnmarshalObject(bz []byte, ptr interface{}) error {
	if len(bz) == 0 {
		err := errors.New("TMClient#UnmarshalObject []byte data is empty")
		// fmt.Println("[ERROR] ", err)
		return err
	}
	
	err := client.codec.UnmarshalBinaryBare(bz, ptr)
	
	if err != nil {
		err := errors.New("TMClient#UnmarshalObject : " + err.Error())
		// fmt.Println("[ERROR] ", err)
		return err
	}
	return err
}