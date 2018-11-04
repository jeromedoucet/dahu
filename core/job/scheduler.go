package job

// the scheduler goal is to keep control over with running execution.
// everything is done through channel till execution are ran in goroutine.
// like for notifier, there is a local goroutine which handle execution
// registration and command on it. Its the only process that should change
// inner state of scheduler.

type schedulingRegistration struct {
	id string
	c  chan interface{}
}

var runningJobExecution map[string]chan interface{}
var registerChan chan schedulingRegistration
var unregisterChan chan string
var cancelationChan chan string

func init() {
	runningJobExecution = make(map[string]chan interface{})
	registerChan = make(chan schedulingRegistration)
	unregisterChan = make(chan string)
	cancelationChan = make(chan string)
	go startScheduler()
}

func startScheduler() {
	for {
		select {
		case reg := <-registerChan:
			runningJobExecution[reg.id] = reg.c
		case id := <-unregisterChan:
			delete(runningJobExecution, id)
		case id := <-cancelationChan:
			c, ok := runningJobExecution[id]
			if ok {
				close(c)
				delete(runningJobExecution, id)
			}
		}
	}
}

// registerJobExecution add a jobExecution watcher for
// being able to interract with. Return a chan that is used
// to notify cancelation
func registerJobExecution(jobId, jobExecutionId string) chan interface{} {
	c := make(chan interface{})
	registerChan <- schedulingRegistration{id: jobId + jobExecutionId, c: c}
	return c
}

// unRegisterJobExecution allow to unregister a job execution
// from the inner side. Should be used when a job end.
func unRegisterJobExecution(jobId, jobExecutionId string) {
	unregisterChan <- jobId + jobExecutionId
}

// AskForCancelation is called to cancel a jobExecution.
// if there is no corresponding running jobExecution,
// the demand is simply ignored.
func AskForCancelation(jobId, jobExecutionId string) {
	cancelationChan <- jobId + jobExecutionId
}
