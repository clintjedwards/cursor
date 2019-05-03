package storage

import (
	"fmt"
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

	err = boltDB.createBuckets("pipelines")
	if err != nil {
		return err
	}

	return nil
}

func (boltDB *boltDB) createBuckets(names ...string) error {

	for _, name := range names {
		err := boltDB.store.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(name))
			if err != nil {
				return fmt.Errorf("could not create bucket: %s; %v", name, err)
			}

			return nil
		})

		if err != nil {
			return err
		}
	}
	return nil
}

func (boltDB *boltDB) GetAll() map[string][]byte {
	return nil
}

func (boltDB *boltDB) Get(key string) ([]byte, error) {
	return nil, nil
}

func (boltDB *boltDB) Add(bucketName string, key, value []byte) error {

	err := boltDB.store.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		err := bucket.Put(key, value)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

func (boltDB *boltDB) Update(key string, newValue []byte) error {
	return nil
}

func (boltDB *boltDB) Delete(key string) error {
	return nil
}
