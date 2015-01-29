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

package local

import (
    "os"
    "strings"

    "github.com/lighthouse/beacon/structs"
)

var Driver = &structs.Driver {
    Name: "local",
    IsApplicable: IsApplicable,
    GetVMs: GetVMS,
}

func IsApplicable() bool {
    return os.Getenv("DOCKER_HOST") != ""
}

func GetVMS() []*structs.VM {
    dockerHost := os.Getenv("DOCKER_HOST")
    dockerHost = strings.Replace(dockerHost, "tcp://", "", 1)
    hostInfo := strings.Split(dockerHost, ":")

    boot2Docker := &structs.VM{
        Name: "boot2docker",
        Address: hostInfo[0],
        Port: hostInfo[1],
        Version: "v1",
        CanAccessDocker: false,
    }
    boot2Docker.CanAccessDocker = boot2Docker.PingDocker()

    if boot2Docker.CanAccessDocker {
        boot2Docker.Version, _ = boot2Docker.GetDockerVersion()
    }

    return []*structs.VM{boot2Docker}
}
