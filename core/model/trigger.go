package model

// use by request that
// will start a job run.
type RunRequest struct {
	OpenWs bool `json:"openWs"`
}

type JobRun struct {
	RunId         string `json:"runId"`
	ContainerName string `json:"containerName"`
}
