package types

type ViewType uint8

const (
	GetOne  = ViewType(1)
	GetMany = ViewType(2)
)

const (
	MethodPutHeartbeat = "daemon/heartbeat"
)

const (
	PathMember         = "$sys/member"
)



type TxMsg struct {
	Method string
	Path   string
	Key    []byte
	Data   []byte
}

func NewTxMsg(method string, path string, key string, data []byte) *TxMsg {
	msg := TxMsg{Method: method, Path: path, Key: []byte(key), Data: data}
	return &msg
}

func EncodeTx(method string, path string, key string, data []byte) ([]byte, error) {
	msg := TxMsg{Method: method, Path: path, Key: []byte(key), Data: data}
	return BasicCdc.MarshalBinaryBare(msg)
}

func EncodeTxMsg(msg *TxMsg) ([]byte, error) {
	return BasicCdc.MarshalBinaryBare(msg)
}

func DecodeTxMsg(msgBytes []byte) (*TxMsg, error) {
	msg := &TxMsg{}
	err := BasicCdc.UnmarshalBinaryBare(msgBytes, msg)
	return msg, err
}

type ViewMsg struct {
	Type   ViewType
	Path   string
	Start  []byte
	End    []byte
}

func NewViewMsgOne(path string, key string) *ViewMsg {
	msg := ViewMsg{Type:GetOne, Path: path, Start: []byte(key)}
	return &msg
}

func NewViewMsgMany(path string, start string, end string) *ViewMsg {
	var endBytes []byte
	if len(end) > 0{
		endBytes = []byte(end)
	}
	msg := ViewMsg{Type:GetMany, Path: path, Start: []byte(start), End:endBytes}
	return &msg
}

func EncodeViewMsg(msg *ViewMsg) ([]byte, error) {
	return BasicCdc.MarshalBinaryBare(msg)
}

func DecodeViewMsg(msgBytes []byte) (*ViewMsg, error) {
	msg := &ViewMsg{}
	err := BasicCdc.UnmarshalBinaryBare(msgBytes, msg)
	return msg, err
}
