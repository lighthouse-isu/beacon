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

package ocean

import (
    "fmt"
    "sync"
    "time"
    "net/http"
    "io/ioutil"
    "encoding/json"

    "github.com/lighthouse/beacon/structs"
)

var oceanToken = ""

var Driver = &structs.Driver {
    Name: "ocean",
    IsApplicable: IsApplicable,
    GetVMs: GetProjectVMs,
}

type DropletNetwork struct {
    IP string `json:"ip_address,omitempty"`
    Gateway string `json:"gateway,omitempty"`
    Type string `json:"type,omitempty"`
}

type DropletNetworks struct {
    V4 []*DropletNetwork `json:"v4,omitempty"`
    V6 []*DropletNetwork `json:"v6,omitempty"`
}

type Droplet struct {
    Name string `json:"name,omitempty"`
    Networks DropletNetworks `json:"networks,omitempty"`
}

type DropletList struct {
    Droplets []*Droplet
}

func GetOceanToken() string {
    if oceanToken == "" {
        request, _ := http.NewRequest("GET",
            "http://169.254.169.254/metadata/v1/user-data", nil)

        client := http.Client{}
        resp, err := client.Do(request)
        defer resp.Body.Close()
        if err != nil {
            return ""
        }

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            return ""
        }

        oceanToken = string(body)
    }

    return oceanToken
}

func IsApplicable() bool {
    request, _ := http.NewRequest("GET",
        "http://169.254.169.254/metadata/v1/", nil)

    client := http.Client {
        Timeout: time.Duration(2 * time.Second),
    }
    resp, err := client.Do(request)

    return err == nil && resp.StatusCode == 200
}

func GetProjectVMs() []*structs.VM {
    request, _ := http.NewRequest("GET",
        "https://api.digitalocean.com/v2/droplets", nil)

    request.Header.Add("Authorization",
        fmt.Sprintf("Bearer %s", GetOceanToken()))

    var discoveredVMs []*structs.VM

    client := http.Client{}
    resp, err := client.Do(request)
    if err != nil {
        return discoveredVMs
    }
    body, err := ioutil.ReadAll(resp.Body)
    defer resp.Body.Close()
    if err != nil {
        return discoveredVMs
    }

    var dropletList DropletList
    err = json.Unmarshal(body, &dropletList)
    if err != nil {
        return discoveredVMs
    }

    for _, droplet := range dropletList.Droplets {
        for _, network := range droplet.Networks.V4 {
            if network.Type == "private" {
                discoveredVMs = append(discoveredVMs, &structs.VM{
                    Name: droplet.Name,
                    Address: network.IP,
                    Port: "2375",
                    Version: "v1",
                    CanAccessDocker: false,
                })
            }
        }
    }

    var wg sync.WaitGroup
    for _, vm := range discoveredVMs {
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
