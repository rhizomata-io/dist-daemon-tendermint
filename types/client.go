package types


// DataHandler data handler for GetMany
type KVHandler func(key []byte, value []byte) bool


type Client interface {
	BroadcastTxSync(msg *TxMsg) error
	BroadcastTxAsync(msg *TxMsg) error
	Commit() error
	Query(msg *ViewMsg) ([]byte, error)
	GetObject(msg *ViewMsg, obj interface{}) (err error)
	GetMany(msg *ViewMsg, handler KVHandler) (err error)
	UnmarshalObject(bz []byte, ptr interface{}) error
	MarshalObject(ptr interface{}) ([]byte, error)
}
