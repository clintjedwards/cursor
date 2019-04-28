package master

import (
	"net"

	"github.com/boltdb/bolt"
	"github.com/clintjedwards/cursor/config"
	"github.com/clintjedwards/cursor/storage"
	"github.com/clintjedwards/cursor/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type CursorMaster struct {
	storage *bolt.DB
	config  *config.Config
}

func initCursorMaster() *CursorMaster {
	cursorMaster := CursorMaster{}

	config, err := config.FromEnv()
	if err != nil {
		utils.StructuredLog(utils.LogLevelFatal, "failed to get config", err)
	}

	storage, err := storage.NewBoltDB(config.Database.Path)
	if err != nil {
		utils.StructuredLog(utils.LogLevelFatal, "failed init storage", err)
	}

	cursorMaster.config = config
	cursorMaster.storage = storage

	return &cursorMaster
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

	//initCursorMasterobject
	//RegisterRoutes

	utils.StructuredLog(utils.LogLevelInfo,
		"started cursor master",
		map[string]string{"server_url": config.ServerURL})
	grpcServer.Serve(listen)
}
