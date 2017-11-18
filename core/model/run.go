package model

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync"
	"time"
)

// run params
type ProcessParams struct {
	Id           string
	Image        string
	Env          map[string]string
	OutputWriter io.Writer
	TimeOut      time.Duration
}

func (p *ProcessParams) ContainerName() string {
	return fmt.Sprintf("dahu-run-%s", p.Id)
}

// initiate a new Run with given params
func NewProcess(p ProcessParams) *Process {
	r := new(Process)
	r.params = p
	r.cmdArg = formatProcessParams(r.params)
	r.done = make(chan interface{})
	r.m = &sync.Mutex{}
	r.status = CREATED
	return r
}

// format a docker run command
// with some constants ('run')
// and thanks to run params (env and image name for instance).
func formatProcessParams(p ProcessParams) []string {
	args := []string{"run", "--name", p.ContainerName()}
	if len(p.Env) != 0 {
		buf := make([]string, 0)
		for key, value := range p.Env {
			buf = append(buf, "--env")
			buf = append(buf, key+"="+value)
		}
		args = append(args, buf...)
	}
	return append(args, p.Image)
}

// A container process
// for example, this could be a git clone,
// a run of a test set or even a deployment
type Process struct {
	params ProcessParams
	cmdArg []string
	status RunStatus // must be accessed through thread-safe functions Status() and setStatus()
	done   chan interface{}
	m      *sync.Mutex
	cmd    *exec.Cmd
}

// start the command.
// return error immediatly if the
// status of the command has another
// value than CREATED.
// this function is thread-safe and
// non blocking.
func (r *Process) Start(ctx context.Context) error {
	// todo return JobRun ?
	r.m.Lock()
	defer r.m.Unlock()
	if r.status == CREATED {
		timeOut := r.params.TimeOut
		// no 0 or negative timeOut allowed
		if timeOut <= time.Duration(0) {
			timeOut = defaultTimeOut
		}
		c, cancel := context.WithTimeout(ctx, timeOut)
		r.cmd = exec.CommandContext(c, "docker", r.cmdArg...)
		// todo check the output ! and test it
		// think on Stderr ?
		r.cmd.Stdout = r.params.OutputWriter
		r.cmd.Stderr = r.params.OutputWriter

		log.Printf("INFO >> run now the command : %+v", r.cmd.Args)
		r.cmd.Start()
		r.status = RUNNING
		go func() {
			defer close(r.done)
			err := r.cmd.Wait() // todo handle error here
			cancel()
			r.m.Lock()
			defer r.m.Unlock()
			log.Printf("INFO >> teminate the comand whith error : %+v", err)
			if r.cmd.ProcessState.Success() {
				r.status = SUCCESS
			} else if r.status != CANCELED { // if the command has already been canceled, must not change status
				r.status = FAILURE
				if c.Err() == context.DeadlineExceeded {
					r.params.OutputWriter.Write([]byte("Time out"))
				}
			}
		}()
		return nil
	} else {
		msg := "ERROR >> model.Run.Start try to Start a Run more than one time"
		log.Print(msg)
		return errors.New(msg)
	}
}

// return the status of the
// command. This function is
// thread-safe.
func (r *Process) Status() RunStatus {
	r.m.Lock()
	defer r.m.Unlock()
	return r.status
}

// change the status of the
// command.
// This function is thread-safe,
// permitting to make some change
// without knowing the internal state
// of the command.
// It is usefull for some action that
// must be done internally too.
func (r *Process) setStatus(s RunStatus) {
	r.m.Lock()
	defer r.m.Unlock()
	r.status = s
}

// return a channel that
// allow to check if the
// command has finished.
// It is done when the return
// channel is closed.
func (r *Process) Done() chan interface{} {
	return r.done
}

// cancel a running command.
// if trying to cancel a non-started
// or a finished command, it will
// return an error.
// return an error if the command
// process termination encounter some issues.
//
// this function is thread-safe
func (r *Process) Cancel() error {
	r.m.Lock()
	defer r.m.Unlock()
	if r.status != RUNNING {
		return errors.New("can only cancel a RUNNING Run")
	}
	err := r.cmd.Process.Kill()
	// todo make something with that error
	// todo test this err
	if err == nil {
		r.status = CANCELED
	}
	return err
}
