package sdk

import (
	cursorPlugin "github.com/clintjedwards/cursor/plugin"
	proto "github.com/clintjedwards/cursor/plugin/proto"
	"github.com/hashicorp/go-plugin"
)

// Task is the smallest unit of work. A task accomplishes some hyper specific thing.
// Many tasks together make up a single pipeline
type Task struct {
	Name        string
	Description string
	Handler     func() error
	Children    []string // Task IDs that should be run after this task completes
}

// Pipeline is used to define overall elements about the pipeline
type Pipeline struct {
	Name     string // Pipeline name TODO: unsure if this matters yet
	RootTask string // ID of task that pipeline will start with by default
	TaskMap  map[string]Task
}

// GetPipelineInfo is called upon compiling the pipeline so that cursor knowns the structure of the
// tasks.
// It returns a task info map that describes what kind of tasks exist in the pipeline and a
// default task used as the root task
func (pipeline *Pipeline) GetPipelineInfo(request *proto.GetPipelineInfoRequest) (*proto.GetPipelineInfoResponse, error) {
	pipelineInfo := proto.GetPipelineInfoResponse{
		RootTask: pipeline.RootTask,
		Tasks:    map[string]*proto.TaskInfo{},
	}

	for taskID, task := range pipeline.TaskMap {
		pipelineInfo.Tasks[taskID] = &proto.TaskInfo{
			Name:        task.Name,
			Description: task.Description,
			Children:    task.Children,
		}
	}

	return &pipelineInfo, nil
}

// ExecuteTask finds the appropriate task from the global task map and runs it
// Returns an error if task does not exist in map
func (pipeline *Pipeline) ExecuteTask(taskRequest *proto.ExecuteTaskRequest) (*proto.ExecuteTaskResponse, error) {

	if _, ok := pipeline.TaskMap[taskRequest.Id]; !ok {
		return &proto.ExecuteTaskResponse{
			Status: proto.ExecuteTaskResponse_FAILED,
		}, nil
	}

	task := pipeline.TaskMap[taskRequest.Id]

	err := task.Handler()
	if err != nil {
		return &proto.ExecuteTaskResponse{
			Status: proto.ExecuteTaskResponse_FAILED,
		}, nil
	}

	return &proto.ExecuteTaskResponse{
		Status: proto.ExecuteTaskResponse_SUCCESS,
	}, nil
}

// Serve starts the plugin gRPC server and handles requests
// Should be called last in user pipeline definitions
func Serve(pipelineName, rootTask string, taskMap map[string]Task) {

	pipeline := Pipeline{
		Name:     pipelineName,
		RootTask: rootTask,
		TaskMap:  taskMap,
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: cursorPlugin.Handshake,
		Plugins: map[string]plugin.Plugin{
			// the key here doesn't seem to matter
			"cursor-sdk": &cursorPlugin.CursorPlugin{Impl: &pipeline},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})

}
