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

package structs

import (
    "fmt"
    "time"
    "net/http"
    "io/ioutil"
    "encoding/json"
)

type VM struct {
    Name string `json:"name,omitempty"`
    Address string `json:"address,omitempty"`
    Port string `json:"port,omitempty"`
    Version string `json:"version,omitempty"`
    CanAccessDocker bool `json:"canAccessDocker"`
}

func (this *VM) PingDocker() bool {
    pingAddress := fmt.Sprintf("http://%s:%s/%s/_ping",
        this.Address, this.Port, this.Version)

    client := &http.Client{
        Timeout: time.Duration(2)*time.Second,
    }
    response, err := client.Get(pingAddress)

    return err == nil && response.StatusCode == 200
}

func (this *VM) GetDockerVersion() (string, error) {
    dockerAddress := fmt.Sprintf("http://%s:%s/%s/version",
        this.Address, this.Port, this.Version)

    req, err := http.NewRequest("GET", dockerAddress, nil)
    if err != nil {
        return "v1", err
    }

    resp, err := http.DefaultClient.Do(req)
    defer resp.Body.Close()

    if err != nil {
        return "v1", err
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "v1", err
    }

    var api struct {
        ApiVersion string
    }

    err = json.Unmarshal(body, &api)
    if err != nil {
        return "v1", err
    }
    
    return fmt.Sprintf("v%s", api.ApiVersion), nil
}
