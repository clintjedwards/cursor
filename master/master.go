package master

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/clintjedwards/cursor/api"
	"github.com/clintjedwards/cursor/config"
	"github.com/clintjedwards/cursor/storage"
	"github.com/clintjedwards/cursor/utils"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
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

// CreatePipeline registers a new pipeline
func (master *CursorMaster) CreatePipeline(context context.Context, request *api.CreatePipelineRequest) (*api.CreatePipelineResponse, error) {

	newPipeline := api.Pipeline{
		Id:          string(utils.GenerateRandString(master.config.Master.IDLength)),
		Name:        request.Name,
		Description: request.Description,
		Created:     time.Now().Unix(),
		Modified:    time.Now().Unix(),
		GitRepo: &api.GitRepo{
			Url:    request.GitRepo.Url,
			Branch: request.GitRepo.Branch,
		},
	}

	if newPipeline.GitRepo.Url == "" {
		return &api.CreatePipelineResponse{}, status.Error(codes.FailedPrecondition, "git url required")
	}

	if newPipeline.Name == "" {
		return &api.CreatePipelineResponse{}, status.Error(codes.FailedPrecondition, "name required")
	}

	if newPipeline.GitRepo.Branch == "" {
		newPipeline.GitRepo.Branch = "master"
	}

	protoNewPipeline, err := proto.Marshal(&newPipeline)
	if err != nil {
		return &api.CreatePipelineResponse{}, status.Error(codes.Internal, "failed to marshal proto")
	}

	err = master.storage.Add(storage.PipelinesBucket, newPipeline.Id, protoNewPipeline)
	if err != nil {
		return &api.CreatePipelineResponse{}, status.Error(codes.Internal, "could not save pipeline")
	}

	utils.StructuredLog(utils.LogLevelInfo, "pipeline created", newPipeline)

	return &api.CreatePipelineResponse{Id: newPipeline.Id}, nil
}

// ListPipelines returns a list of all pipelines on a cursor master
func (master *CursorMaster) ListPipelines(context context.Context, request *api.ListPipelinesRequest) (*api.ListPipelinesResponse, error) {

	rawPipelines, err := master.storage.GetAll(storage.PipelinesBucket)
	if err != nil {
		return &api.ListPipelinesResponse{}, status.Error(codes.Internal, "failed to retrieve pipelines from database")
	}

	pipelines := map[string]*api.Pipeline{}

	for key, value := range rawPipelines {
		pipeline := &api.Pipeline{}
		err := proto.Unmarshal(value, pipeline)
		if err != nil {
			utils.StructuredLog(utils.LogLevelError, "could not unmarshal pipeline", map[string]string{"pipeline_key": key, "err": err.Error()})
			continue
		}
		pipelines[key] = pipeline
	}

	return &api.ListPipelinesResponse{Pipelines: pipelines}, nil
}

// GetPipeline returns a single pipeline by id
func (master *CursorMaster) GetPipeline(context context.Context, request *api.GetPipelineRequest) (*api.GetPipelineResponse, error) {
	return nil, nil
}

// TODO: Remember to remove the plugins here also
func (master *CursorMaster) DeletePipeline(context context.Context, request *api.DeletePipelineRequest) (*api.DeletePipelineResponse, error) {
	return nil, nil
}

// RegisterWorker connects a new worker with a given cursor master
func (master *CursorMaster) RegisterWorker(context context.Context, request *api.RegisterWorkerRequest) (*api.RegisterWorkerResponse, error) {
	return nil, nil
}

// StartServer initializes a GRPC server
func StartServer() {
	config, err := config.FromEnv()
	if err != nil {
		utils.StructuredLog(utils.LogLevelFatal, "failed to get config", err)
	}

	listen, err := net.Listen("tcp", config.Master.URL)
	if err != nil {
		utils.StructuredLog(utils.LogLevelFatal, "could not initialize tcp listener", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	cursorMaster := initCursorMaster()

	api.RegisterCursorMasterServer(grpcServer, cursorMaster)

	utils.StructuredLog(utils.LogLevelInfo,
		"started cursor master",
		map[string]string{"server_url": config.Master.URL})
	grpcServer.Serve(listen)
}
