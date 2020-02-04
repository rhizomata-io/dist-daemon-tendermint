package store

import (
	"errors"
	"fmt"
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
	
	dbm "github.com/tendermint/tm-db"
	"sync"
)

const (
	pathKeySeparator = byte('!')
)

type Path []byte

//ValueOf
func valueOf(keyStr string) Path {
	return Path([]byte(keyStr))
}

func (path Path) String() string {
	return string(path)
}

func (path Path) makeKeyString(key string) []byte {
	return path.makeKey([]byte(key))
}

func (path Path) makeKey(key []byte) []byte {
	newKey := append(path, pathKeySeparator)
	newKey = append(newKey, key...)
	return newKey
}

type Store struct {
	path Path
	db   dbm.DB
}

func NewStore(path string, db dbm.DB) *Store {
	return &Store{path: valueOf(path), db: db}
}

func (store *Store) GetPath() string {
	return store.path.String()
}

func (store *Store) Get(key []byte) ([]byte, error) {
	keyBytes := store.path.makeKey(key)
	return store.db.Get(keyBytes)
}

func (store *Store) Has(key []byte) (bool, error) {
	keyBytes := store.path.makeKey(key)
	return store.db.Has(keyBytes)
}

func (store *Store) Set(key []byte, value []byte) error {
	keyBytes := store.path.makeKey(key)
	return store.db.Set(keyBytes, value)
}

func (store *Store) SetSync(key []byte, value []byte) error {
	keyBytes := store.path.makeKey(key)
	return store.db.SetSync(keyBytes, value)
}

func (store *Store) Delete(key []byte) error {
	keyBytes := store.path.makeKey(key)
	return store.db.Delete(keyBytes)
}

func (store *Store) DeleteSync(key []byte) error {
	keyBytes := store.path.makeKey(key)
	return store.db.DeleteSync(keyBytes)
}

func (store *Store) Iterator(start, end []byte) (dbm.Iterator, error) {
	startBytes := store.path.makeKey(start)
	var endBytes []byte
	if len(end) > 0 {
		endBytes = store.path.makeKey(end)
	}
	
	return store.db.Iterator(startBytes, endBytes)
}


func (store *Store) GetMany(start, end []byte) (arrayBytes[]byte, err error) {
	iterator, err := store.Iterator(start, end)
	
	if err != nil {
		return nil,err
	}
	
	recordsArray := [][]byte{}
	
	for iterator.Valid() {
		recordsArray = append(recordsArray, iterator.Value())
		//fmt.Println(" *** Iterator:", string(iterator.Value()))
		iterator.Next()
	}
	
	arrayBytes, err = types.BasicCdc.MarshalBinaryBare(recordsArray)
	
	return arrayBytes, err
}



type Registry struct {
	sync.Mutex
	pathStores map[string]*Store
	db         dbm.DB
}

func NewRegistry(db dbm.DB) *Registry {
	return &Registry{pathStores: make(map[string]*Store), db: db}
}

func (reg *Registry) RegisterStore(path string) error {
	reg.Lock()
	defer reg.Unlock()
	
	_, ok := reg.pathStores[path]
	if ok {
		return errors.New(fmt.Sprintf("Store [%s] is already registered.", path))
	}
	
	store := &Store{path: valueOf(path), db: reg.db}
	reg.pathStores[path] = store
	return nil
}

func (reg *Registry) GetStore(path string) (*Store, error) {
	reg.Lock()
	defer reg.Unlock()
	
	store, ok := reg.pathStores[path]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Cannot find Store [%s].", path))
	}
	
	return store, nil
}

func (reg *Registry) GetOrMakeStore(path string) *Store {
	reg.Lock()
	defer reg.Unlock()
	
	store, ok := reg.pathStores[path]
	if !ok {
		store = &Store{path: valueOf(path), db: reg.db}
		reg.pathStores[path] = store
	}
	return store
}
