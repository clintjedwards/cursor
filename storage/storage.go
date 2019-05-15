package storage

import (
	"fmt"

	"github.com/clintjedwards/cursor/config"
)

// Bucket represents the name of a section of key/value pairs
type Bucket string

const (
	// PipelinesBucket represents the pipelines bucket
	PipelinesBucket Bucket = "pipelines"
)

// EngineType type represents the different possible storage engines available
type EngineType string

const (
	// StorageEngineBoltDB represents a boltDB storage engine.
	// A file based key-value store.(https://github.com/boltdb/bolt)
	StorageEngineBoltDB EngineType = "boltdb"
)

// Engine represents backend storage implementations where items can be persisted
type Engine interface {
	Init(config *config.Config) error
	GetAll(bucket Bucket) (map[string][]byte, error)
	Get(bucket Bucket, key string) ([]byte, error)
	Add(bucket Bucket, key string, value []byte) error
	Update(bucket Bucket, key string, newValue []byte) error
	Delete(bucket Bucket, key string) error
}

// InitStorage creates a storage object with the appropriate engine
func InitStorage(engineType EngineType) (Engine, error) {

	switch engineType {
	case StorageEngineBoltDB:
		config, err := config.FromEnv()
		if err != nil {
			return nil, err
		}

		boltDBStorageEngine := boltDB{}
		err = boltDBStorageEngine.Init(config)
		if err != nil {
			return nil, err
		}

		return &boltDBStorageEngine, nil
	default:
		return nil, fmt.Errorf("storage backend not implemented: %s", engineType)
	}
}
