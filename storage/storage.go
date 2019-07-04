package storage

import (
	"fmt"

	"github.com/clintjedwards/cursor/api"
	"github.com/clintjedwards/cursor/config"
)

// Bucket represents the name of a section of key/value pairs
// usually a grouping of some sort
// ex. A key/value pair of userid-userdata would belong in the users bucket
type Bucket string

const (
	// PipelinesBucket represents the container in which pipelines are managed
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
	GetAllPipelines(user string) (map[string]*api.Pipeline, error)
	GetPipelines(user, id string) (*api.Pipeline, error)
	AddPipelines(user, id string, pipeline *api.Pipeline) error
	UpdatePipelines(user, id string, pipeline *api.Pipeline) error
	DeletePipelines(user, id string) error
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
