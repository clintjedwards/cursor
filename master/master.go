package master

import (
	"context"
	"net"

	"github.com/clintjedwards/cursor/api"
	"github.com/clintjedwards/cursor/config"
	"github.com/clintjedwards/cursor/storage"
	"github.com/clintjedwards/cursor/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// CursorMaster represents a cursor master server
type CursorMaster struct {
	storage storage.Engine
	config  *config.Config
}

func initCursorMaster() *CursorMaster {
	cursorMaster := CursorMaster{}

	config, err := config.FromEnv()
	if err != nil {
		utils.StructuredLog(utils.LogLevelFatal, "failed to get config", err)
	}

	storage, err := storage.InitStorage(storage.StorageEngineBoltDB)
	if err != nil {
		utils.StructuredLog(utils.LogLevelFatal, "failed init storage", err)
	}

	cursorMaster.config = config
	cursorMaster.storage = storage

	return &cursorMaster
}

// CreatePipeline registers a new pipeline
func (master *CursorMaster) CreatePipeline(context context.Context, request *api.CreatePipelineRequest) (*api.CreatePipelineResponse, error) {

	// attempt to compile github repo saving the resulting binary to a specific directory called plugins
	// we should only try to compile a specified folder that we determine
	// pipelines all have a unique name
	return nil, nil
}

// StartServer initializes a GRPC server
func StartServer() {
	config, err := config.FromEnv()
	if err != nil {
		utils.StructuredLog(utils.LogLevelFatal, "failed to get config", err)
	}

	listen, err := net.Listen("tcp", config.ServerURL)
	if err != nil {
		utils.StructuredLog(utils.LogLevelFatal, "could not initialize tcp listener", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	cursorMaster := initCursorMaster()

	api.RegisterCursorMasterServer(grpcServer, cursorMaster)

	utils.StructuredLog(utils.LogLevelInfo,
		"started cursor master",
		map[string]string{"server_url": config.ServerURL})
	grpcServer.Serve(listen)
}
