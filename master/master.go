package master

import (
	"log"
	"net"

	"github.com/clintjedwards/cursor/api"
	"github.com/clintjedwards/cursor/config"
	"github.com/clintjedwards/cursor/storage"
	"github.com/clintjedwards/cursor/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

// CursorMaster represents a cursor master server
type CursorMaster struct {
	storage storage.Engine
	config  *config.Config
}

// NewCursorMaster inits a grpc cursor master server
func NewCursorMaster(config *config.Config) *CursorMaster {
	cursorMaster := CursorMaster{}

	storage, err := storage.InitStorage(storage.StorageEngineBoltDB)
	if err != nil {
		utils.StructuredLog(utils.LogLevelFatal, "failed init storage", err)
	}

	cursorMaster.config = config
	cursorMaster.storage = storage

	return &cursorMaster
}

// CreateGRPCServer creates a grpc server with all the proper settings; TLS enabled
func CreateGRPCServer(cursorMaster *CursorMaster) *grpc.Server {

	creds, err := credentials.NewServerTLSFromFile(cursorMaster.config.TLSCertPath, cursorMaster.config.TLSKeyPath)
	if err != nil {
		utils.StructuredLog(utils.LogLevelFatal, "failed to get certificates", err)
	}

	serverOption := grpc.Creds(creds)

	grpcServer := grpc.NewServer(serverOption)

	reflection.Register(grpcServer)
	api.RegisterCursorMasterServer(grpcServer, cursorMaster)

	return grpcServer
}

// InitGRPCService starts a GPRC server
func InitGRPCService(config *config.Config, server *grpc.Server) {

	listen, err := net.Listen("tcp", config.Master.GRPCURL)
	if err != nil {
		utils.StructuredLog(utils.LogLevelFatal, "could not initialize tcp listener", err)
	}

	utils.StructuredLog(utils.LogLevelInfo, "starting cursor master grpc service",
		map[string]string{"url": config.Master.GRPCURL})

	log.Fatal(server.Serve(listen))
}
