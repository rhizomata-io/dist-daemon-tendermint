package types


const (
	KeyDelimiter = byte('/')
)

type Key []byte

//ValueOf
func ValueOf(keyStr string) Key {
	return Key([]byte(keyStr))
}

func AppendKey(key Key, next []byte) Key {
	newKey := append(key, KeyDelimiter)
	newKey = append(newKey , next...)
	return newKey
}

func AppendString(key Key, next string) Key {
	return AppendKey(key, ValueOf(next))
}

func (key Key) String() string{
	return string(key)
}

