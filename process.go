// Copyright 2017 CoreSwitch.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package process

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"golang.org/x/net/context"
)

type Process struct {
	Name       string
	Vrf        string
	Args       []string
	File       string
	ErrLookup  string
	ErrStart   string
	ErrWait    string
	ExitFunc   func()
	State      int
	Cmd        *exec.Cmd
	StartTimer int
	RetryTimer int
}

type ProcessSlice []*Process

const (
	PROCESS_STARTING = iota
	PROCESS_RUNNING
	PROCESS_RETRY
	PROCESS_EXIT_CALLED
	PROCESS_STOP_WAIT
	PROCESS_STOP
)

var (
	ProcessList  = ProcessSlice{}
	ProcessMutex sync.RWMutex
)

var ProcessStateStr = map[int]string{
	PROCESS_STARTING:    "Starting",
	PROCESS_RUNNING:     "Running",
	PROCESS_RETRY:       "Retry",
	PROCESS_EXIT_CALLED: "Exit Called",
	PROCESS_STOP_WAIT:   "Stop Wait",
	PROCESS_STOP:        "Stop",
}

func NewProcess(name string, args ...string) *Process {
	proc := &Process{
		Name:       name,
		Args:       []string(args),
		RetryTimer: 1,
	}
	return proc
}

func ProcessRegister(proc *Process) {
	ProcessMutex.Lock()
	defer ProcessMutex.Unlock()

	ProcessList = append(ProcessList, proc)
	proc.Start()
}

func ProcessUnregister(proc *Process) {
	ProcessMutex.Lock()
	defer ProcessMutex.Unlock()

	proc.Stop()

	procList := ProcessSlice{}
	for _, p := range ProcessList {
		if p != proc {
			procList = append(procList, p)
		}
	}
	ProcessList = procList
}

func ProcessCount() int {
	ProcessMutex.Lock()
	defer ProcessMutex.Unlock()

	return len(ProcessList)
}

func ProcessStart(index int) {
	ProcessMutex.Lock()
	defer ProcessMutex.Unlock()

	if index <= 0 {
		return
	}
	if len(ProcessList) < index {
		return
	}
	proc := ProcessList[index-1]
	if proc == nil {
		return
	}
	proc.Start()
}

func ProcessStop(index int) {
	ProcessMutex.Lock()
	defer ProcessMutex.Unlock()

	if index <= 0 {
		return
	}
	if len(ProcessList) < index {
		return
	}
	proc := ProcessList[index-1]
	if proc == nil {
		return
	}
	proc.Stop()
}

func (proc *Process) Start() {
	if proc.ExitFunc != nil {
		return
	}

	proc.State = PROCESS_STOP
	binary, err := exec.LookPath(proc.Name)
	if err != nil {
		proc.ErrLookup = err.Error()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			proc.State = PROCESS_STARTING
			if proc.File != "" {
				os.OpenFile(proc.File, os.O_RDWR|os.O_CREATE, 0644)
			}

			cmd := exec.CommandContext(ctx, binary, proc.Args...)

			env := os.Environ()
			if proc.Vrf != "" {
				env = append(env, fmt.Sprintf("VRF=%s", proc.Vrf))
				env = append(env, "LD_PRELOAD=/usr/bin/vrf_socket.so")
			}
			cmd.Env = env
			proc.Cmd = cmd

			if proc.StartTimer != 0 {
				time.Sleep(time.Duration(proc.StartTimer) * time.Second)
			}

			// fmt.Println("process:", cmd.Path, cmd.Args)
			err = cmd.Start()
			if err != nil {
				proc.ErrStart = err.Error()
			}

			proc.State = PROCESS_RUNNING
			err = cmd.Wait()
			if err != nil {
				proc.ErrWait = err.Error()
			}

			proc.State = PROCESS_RETRY
			retryTimer := time.NewTimer(time.Duration(proc.RetryTimer) * time.Second)
			select {
			case <-retryTimer.C:
			case <-done:
				retryTimer.Stop()
				return
			}
		}
	}()

	proc.ExitFunc = func() {
		proc.State = PROCESS_EXIT_CALLED
		close(done)
		cancel()
		proc.State = PROCESS_STOP_WAIT
		wg.Wait()
		proc.State = PROCESS_STOP
	}
}

func (proc *Process) Stop() {
	if proc.ExitFunc != nil {
		proc.ExitFunc()
		proc.ExitFunc = nil
	}
}

func ProcessListShow() string {
	str := ""
	for pos, proc := range ProcessList {
		str += fmt.Sprintf("%d %s", pos+1, proc.Name)
		if proc.Vrf != "" {
			str += fmt.Sprintf("@%s", proc.Vrf)
		}
		str += fmt.Sprintf(": %s", ProcessStateStr[proc.State])
		if proc.State == PROCESS_RUNNING && proc.Cmd != nil && proc.Cmd.Process != nil {
			str += fmt.Sprintf(" (pid %d)", proc.Cmd.Process.Pid)
		}
		str += "\n"
		if proc.ErrLookup != "" {
			str += fmt.Sprintf("  Last Lookup Error: %s\n", proc.ErrLookup)
		}
		if proc.ErrStart != "" {
			str += fmt.Sprintf("  Last Start Error: %s\n", proc.ErrStart)
		}
		if proc.ErrWait != "" {
			str += fmt.Sprintf("  Last Wait Error: %s\n", proc.ErrWait)
		}
		str += fmt.Sprintf("  %s\n", proc.Args)
	}
	return str
}
