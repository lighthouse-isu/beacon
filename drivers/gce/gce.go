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

package gce

import (
    "sync"
    "net/http"
    "io/ioutil"

    "github.com/lighthouse/beacon/structs"

    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/cloud/compute/metadata"
    compute "google.golang.org/api/compute/v1"
)


var Driver = &structs.Driver {
    Name: "gce",
    IsApplicable: IsApplicable,
    GetVMs: GetProjectVMs,
}

func IsApplicable() bool {
    return metadata.OnGCE()
}

func GetCurrentProjectID() (string, error) {
    request, _ := http.NewRequest(
        "GET", "http://metadata.google.internal/computeMetadata/v1/project/project-id", nil)

    request.Header.Add("Metadata-Flavor", "Google")

    client := http.Client{}
    response, err := client.Do(request)
    if err != nil {
        return "", err
    }

    projectID, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return "", err
    }
    response.Body.Close()

    return string(projectID), nil
}

func GetProjectVMs() []*structs.VM {
    client := &http.Client{
        Transport: &oauth2.Transport{
            Source: google.ComputeTokenSource(""),
        },
    }
    computeClient, _ := compute.New(client)

    projectID, _ := GetCurrentProjectID()

    zones, _ := computeClient.Instances.AggregatedList(projectID).Do()

    var discoveredVMs []*structs.VM

    for _, zone := range zones.Items {
        for _, instance := range zone.Instances {
            // For future reference we need figure out which network interface
            // to use instead of deafulting to the first one.
            network := instance.NetworkInterfaces[0]

            discoveredVMs = append(discoveredVMs, &structs.VM{
                Name: instance.Name,
                Address: network.NetworkIP,
                Port: "2375",
                Version: "v1",
                CanAccessDocker: false,
            })
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
