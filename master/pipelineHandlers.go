package master

import (
	"context"
	"log"
	"time"

	"github.com/clintjedwards/cursor/api"
	"github.com/clintjedwards/cursor/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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

	err := master.cloneRepository(newPipeline.Id, newPipeline.GitRepo)
	if err != nil {
		log.Println(err)
	}

	err = master.buildPlugin(newPipeline.Id)
	if err != nil {
		log.Println(err)
	}

	err = master.storage.AddPipeline(newPipeline.Id, &newPipeline)
	if err != nil {
		return &api.CreatePipelineResponse{}, status.Error(codes.Internal, "could not save pipeline when attempting to create new pipeline")
	}

	utils.StructuredLog(utils.LogLevelInfo, "pipeline created", newPipeline)

	return &api.CreatePipelineResponse{Id: newPipeline.Id}, nil
}

// ListPipelines returns a list of all pipelines on a cursor master
func (master *CursorMaster) ListPipelines(context context.Context, request *api.ListPipelinesRequest) (*api.ListPipelinesResponse, error) {

	pipelines, err := master.storage.GetAllPipelines()
	if err != nil {
		utils.StructuredLog(utils.LogLevelError, "could not retrieve pipelines from database while attempting to list all", err)
		return &api.ListPipelinesResponse{}, status.Error(codes.Internal, "failed to retrieve pipelines from database")
	}

	return &api.ListPipelinesResponse{Pipelines: pipelines}, nil
}

// GetPipeline returns a single pipeline by id
func (master *CursorMaster) GetPipeline(context context.Context, request *api.GetPipelineRequest) (*api.GetPipelineResponse, error) {
	if request.Id == "" {
		return &api.GetPipelineResponse{}, status.Error(codes.FailedPrecondition, "pipeline id required")
	}

	pipeline, err := master.storage.GetPipeline(request.Id)
	if err != nil {
		if err == utils.ErrPipelineNotFound {
			return &api.GetPipelineResponse{}, status.Error(codes.NotFound, "requested pipeline not found")
		}
		return &api.GetPipelineResponse{}, status.Error(codes.Internal, "failed to retrieve formula from database")
	}

	return &api.GetPipelineResponse{Pipeline: pipeline}, nil
}

// DeletePipeline removes a pipeline
func (master *CursorMaster) DeletePipeline(context context.Context, request *api.DeletePipelineRequest) (*api.DeletePipelineResponse, error) {
	if request.Id == "" {
		return &api.DeletePipelineResponse{}, status.Error(codes.FailedPrecondition, "pipeline id required")
	}

	err := master.storage.DeletePipeline(request.Id)
	if err != nil {
		if err == utils.ErrPipelineNotFound {
			return &api.DeletePipelineResponse{}, status.Error(codes.NotFound, "could not delete pipeline; key not found")
		}
		utils.StructuredLog(utils.LogLevelError, "could not delete pipeline", err)
		return &api.DeletePipelineResponse{}, status.Error(codes.Internal, "could not delete pipeline")
	}

	utils.StructuredLog(utils.LogLevelInfo, "formula deleted", request.Id)

	return &api.DeletePipelineResponse{}, nil
}

// RegisterWorker connects a new worker with a given cursor master
func (master *CursorMaster) RegisterWorker(context context.Context, request *api.RegisterWorkerRequest) (*api.RegisterWorkerResponse, error) {
	return nil, nil
}
