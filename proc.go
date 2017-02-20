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

package proc

import (
	"sync"
)

type Proc struct {
	Name string
	Vrf  string
	Args []string
	File string
}

type ProcList map[string]*Proc

var (
	ProcVrfMap = map[string]*ProcList{}
	ProcMutex  sync.RWMutex
)

func NewProc(name string, vrf string, args ...string) *Proc {
	proc := &Proc{
		Name: name,
		Vrf:  vrf,
		Args: []string(args),
	}
	return proc
}

func ProcRegister(proc *Proc) {
	ProcMutex.Lock()
	defer ProcMutex.Unlock()
}

func ProcUnregister(proc *Proc) {
	ProcMutex.Lock()
	defer ProcMutex.Unlock()
}

func ProcUnregisterByCommand(name string) {
	ProcMutex.Lock()
	defer ProcMutex.Unlock()
}

func ProcUnregisterByVrf() {
	ProcMutex.Lock()
	defer ProcMutex.Unlock()
}

func (proc *Proc) Start() {
}

func (proc *Proc) Stop() {
}
