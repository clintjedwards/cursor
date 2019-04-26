package models

type JobStatus string

const (
	JobStatusUnknown JobStatus = "unknown"
	JobStatusError   JobStatus = "failed"
	JobStatusSuccess JobStatus = "success"
	JobStatusRunning JobStatus = "running"
	JobStatusWaiting JobStatus = "waiting"
)

type PipelineStatus string

const (
	PipelineStatusRunning PipelineStatus = "running"
	PipelineStatusReady   PipelineStatus = "ready"
	PipelineStatusError   PipelineStatus = "error"
)

// Job represents a single instnce of a
type Job struct {
	ID          int
	Name        string
	Description string
	Status      JobStatus
	DependsOn   []*Job
}

// Pipeline represents a single instance of a pipeline
type Pipeline struct {
	ID            int
	Name          string
	Description   string
	Created       int
	Modified      int
	LastRunID     int
	LastRunDate   int // Complete time of last run
	LastRunAuthor string
}

// PipelineRun represents a single run of a single pipeline
type PipelineRun struct {
	ID         int
	PipelineID int
	LastJob    int
	Start      int
	End        int
	Status     PipelineStatus
}

// GitRepo represents a single git repository
type GitRepo struct {
	URL      string
	Branch   string
	Branches []string
}
