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

package drivers


import (
    "flag"

    "github.com/lighthouse/beacon/structs"

    "github.com/lighthouse/beacon/drivers/gce"
    "github.com/lighthouse/beacon/drivers/local"
    "github.com/lighthouse/beacon/drivers/unknown"

)

var Preferred = flag.String("driver", "", "Specified driver to use")

var Defaults = []*structs.Driver{
    gce.Driver,
    local.Driver,
}


func Decide() *structs.Driver {
    if *Preferred != "" {
        return Find(*Preferred, Defaults)
    }
    return Guess(Defaults)
}

func Find(preferred string, drivers []*structs.Driver) *structs.Driver {
    for _, driver := range drivers {
        if driver.Name == preferred {
            return driver
        }
    }
    return unknown.Driver
}

func Guess(drivers []*structs.Driver) *structs.Driver {
    for _, driver := range drivers {
        if driver.IsApplicable() {
            return driver
        }
    }
    return unknown.Driver
}
