package master

import (
	"context"
	"fmt"
	"time"

	"github.com/clintjedwards/cursor/api"
	"github.com/clintjedwards/cursor/plugin"
	"github.com/clintjedwards/cursor/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreatePipeline registers a new pipeline
func (master *CursorMaster) CreatePipeline(context context.Context, request *api.CreatePipelineRequest) (*api.CreatePipelineResponse, error) {

	// Validate user input
	if request.RepositoryUrl == "" {
		return &api.CreatePipelineResponse{}, status.Error(codes.FailedPrecondition, "repository url required")
	}

	if request.Name == "" {
		return &api.CreatePipelineResponse{}, status.Error(codes.FailedPrecondition, "name required")
	}

	// Generate Pipeline unique ID
	pipelineID := string(utils.GenerateRandString(master.config.Master.IDLength))

	err := master.getRepository(pipelineID, request.RepositoryUrl)
	if err != nil {
		return &api.CreatePipelineResponse{}, status.Error(codes.Internal, fmt.Sprintf("could not get repository: %s", err))
	}

	err = master.buildPlugin(pipelineID)
	if err != nil {
		return &api.CreatePipelineResponse{}, status.Error(codes.Internal, fmt.Sprintf("could not build plugin: %s", err))
	}

	pipelineInfo, err := master.getPipelineInfo(pipelineID)
	if err != nil {
		return &api.CreatePipelineResponse{}, status.Error(codes.Internal, fmt.Sprintf("could not get pipeline info from plugin: %s", err))
	}

	newPipeline := api.Pipeline{
		Id:            pipelineID,
		Name:          request.Name,
		Description:   request.Description,
		RepositoryUrl: request.RepositoryUrl,
		Created:       time.Now().Unix(),
		Modified:      time.Now().Unix(),
		LastCompiled:  time.Now().Unix(),
		RootTaskId:    pipelineInfo.RootTask,
		Tasks:         map[string]*api.Task{},
	}

	// Convert plugin taskInfo proto to api task proto
	for taskID, task := range pipelineInfo.Tasks {
		newPipeline.Tasks[taskID] = &api.Task{
			Name:        task.Name,
			Description: task.Description,
			Children:    task.Children,
		}
	}

	err = master.storage.AddPipeline(newPipeline.Id, &newPipeline)
	if err != nil {
		return &api.CreatePipelineResponse{}, status.Error(codes.Internal, "could not save pipeline when attempting to create new pipeline")
	}

	master.pluginMapMutex.Lock()
	defer master.pluginMapMutex.Unlock()
	master.pluginMap[newPipeline.Id] = &plugin.CursorPlugin{}

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
		if err == utils.ErrEntityNotFound {
			return &api.GetPipelineResponse{}, status.Error(codes.NotFound, "requested pipeline not found")
		}
		return &api.GetPipelineResponse{}, status.Error(codes.Internal, "failed to retrieve formula from database")
	}

	return &api.GetPipelineResponse{Pipeline: pipeline}, nil
}

// RunPipeline triggers a specific pipeline to run
func (master *CursorMaster) RunPipeline(context context.Context, request *api.RunPipelineRequest) (*api.RunPipelineResponse, error) {
	if request.Id == "" {
		return &api.RunPipelineResponse{}, status.Error(codes.FailedPrecondition, "pipeline id required")
	}

	taskID := request.TaskId
	if taskID == "" {
		pipeline, err := master.storage.GetPipeline(request.Id)
		if err != nil {
			if err == utils.ErrEntityNotFound {
				return &api.RunPipelineResponse{}, status.Error(codes.NotFound, "requested pipeline not found")
			}
			return &api.RunPipelineResponse{}, status.Error(codes.Internal, "failed to retrieve formula from database")
		}

		taskID = pipeline.RootTaskId
	}

	err := master.runTasks(request.Id, taskID)
	if err != nil {
		return &api.RunPipelineResponse{}, status.Error(codes.Internal, "failed to run pipeline")
	}

	utils.StructuredLog(utils.LogLevelInfo, "pipeline queued to run", request.Id)

	return &api.RunPipelineResponse{}, nil
}

// DeletePipeline removes a pipeline
func (master *CursorMaster) DeletePipeline(context context.Context, request *api.DeletePipelineRequest) (*api.DeletePipelineResponse, error) {
	if request.Id == "" {
		return &api.DeletePipelineResponse{}, status.Error(codes.FailedPrecondition, "pipeline id required")
	}

	err := master.storage.DeletePipeline(request.Id)
	if err != nil {
		if err == utils.ErrEntityNotFound {
			return &api.DeletePipelineResponse{}, status.Error(codes.NotFound, "could not delete pipeline; key not found")
		}
		utils.StructuredLog(utils.LogLevelError, "could not delete pipeline", err)
		return &api.DeletePipelineResponse{}, status.Error(codes.Internal, "could not delete pipeline")
	}

	master.pluginMapMutex.Lock()
	defer master.pluginMapMutex.Unlock()
	delete(master.pluginMap, request.Id)

	utils.StructuredLog(utils.LogLevelInfo, "formula deleted", request.Id)

	return &api.DeletePipelineResponse{}, nil
}

// GetPipelineRun returns a single pipeline by id
func (master *CursorMaster) GetPipelineRun(context context.Context, request *api.GetPipelineRunRequest) (*api.GetPipelineRunResponse, error) {
	return &api.GetPipelineRunResponse{}, nil
}
