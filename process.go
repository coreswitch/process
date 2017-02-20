// Copyright 2017 CoreSwitch
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
	"sync"
)

type Process struct {
	Name string
	Vrf  string
	Args []string
	File string
}

type ProcessList map[string]*Process

var (
	ProcessVrfMap = map[string]*ProcessList{}
	ProcessMutex  sync.RWMutex
)

func NewProcess(name string, vrf string, args ...string) *Process {
	proc := &Process{
		Name: name,
		Vrf:  vrf,
		Args: []string(args),
	}
	return proc
}

func ProcessRegister(proc *Process) {
	ProcessMutex.Lock()
	defer ProcessMutex.Unlock()
}

func ProcessUnregister(proc *Process) {
	ProcessMutex.Lock()
	defer ProcessMutex.Unlock()
}

func ProcessUnregisterByCommand(name string) {
	ProcessMutex.Lock()
	defer ProcessMutex.Unlock()
}

func ProcessUnregisterByVrf() {
	ProcessMutex.Lock()
	defer ProcessMutex.Unlock()
}

func (proc *Process) Start() {
}

func (proc *Process) Stop() {
}
