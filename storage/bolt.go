package storage

import (
	"time"

	"github.com/boltdb/bolt"
	"github.com/clintjedwards/cursor/config"
)

type boltDB struct {
	filePath string
	store    *bolt.DB
}

func (boltDB *boltDB) Init(config *config.Config) error {

	db, err := bolt.Open(config.Database.Path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	boltDB.store = db

	return nil
}

func (boltDB *boltDB) GetAll() map[string][]byte {
	return nil
}

func (boltDB *boltDB) Get(key string) ([]byte, error) {
	return nil, nil
}

func (boltDB *boltDB) Add(value []byte) error {
	return nil
}

func (boltDB *boltDB) Update(key string, newValue []byte) error {
	return nil
}

func (boltDB *boltDB) Delete(key string) error {
	return nil
}
