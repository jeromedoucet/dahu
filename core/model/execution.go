package model

type EventType string

const (
	JobStart     EventType = "job-start"
	StepStart    EventType = "step-start"
	StepFailed   EventType = "step-failed"
	StepCanceled EventType = "step-canceled"
	StepSucceed  EventType = "step-succeed"
	JobFailed    EventType = "job-failed"
	JobCanceled  EventType = "job-canceled"
	JobSucceed   EventType = "job-succeed"
	NewLog       EventType = "new-log"
)

// Event contains updated information on a current execution.
// it could be job start, new logs, step start, step failed, step succeed,
// job succeed, job failed, ...
type Event struct {
	Type        EventType `json:"type"`
	ExecutionId string    `json:"execution-id"`
	Value       string    `json:"value"`
}
