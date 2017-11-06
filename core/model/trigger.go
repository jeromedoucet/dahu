package model

// use by request that
// will start a job run.
type Trigger struct {
	OpenWs bool `json:"openWs"`
}

type TriggerResponse struct {
	RunId         string `json:"runId"`
	ContainerName string `json:"containerName"`
}
