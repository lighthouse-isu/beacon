// Copyright 2014 Caleb Brose, Chris Fogerty, Rob Sheehy, Zach Taylor, Nick Miller
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

package config

import (
    "os"
    "sync"
    "io/ioutil"
    "encoding/json"

    "github.com/lighthouse/beacon/structs"
)

var Driver = &structs.Driver {
    Name: "config",
    IsApplicable: IsApplicable,
    GetVMs: GetVMs,
}

func IsApplicable() bool {
    _, err := os.Stat("config.json")
    return !os.IsNotExist(err)
}

func GetVMs() []*structs.VM {
    file, e := ioutil.ReadFile("config.json")
    var discoveredVMs []*structs.VM

    if e != nil {
        return discoveredVMs
    }

    var vmIps []string
    json.Unmarshal(file, &vmIps)

    var wg sync.WaitGroup
    for _, vmIp := range vmIps {
        vm := &structs.VM{
            Name: vmIp,
            Address: vmIp,
            Port: "2375",
            Version: "v1",
            CanAccessDocker: false,
        }

        discoveredVMs = append(discoveredVMs, vm)

        wg.Add(1)
        go func(vm *structs.VM) {
            defer wg.Done()
            vm.CanAccessDocker = vm.PingDocker()

            if vm.CanAccessDocker {
                vm.Version, _ = vm.GetDockerVersion()
            }
        }(vm)
    }
    wg.Wait()
    return discoveredVMs
}
