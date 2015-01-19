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

package main

import (
    "fmt"
    "flag"
    "log"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "strings"

    "github.com/mgutz/ansi"

    "github.com/zenazn/goji/web"
    "github.com/zenazn/goji/graceful"
    "github.com/zenazn/goji/web/middleware"

    "github.com/lighthouse/beacon/auth"
    "github.com/lighthouse/beacon/drivers"
    "github.com/lighthouse/beacon/structs"
)


var err error

var pemFile = flag.String("pem", "", "Path to Cert file")
var keyFile = flag.String("key", "", "Path to Key file")
var address = flag.String("h", "127.0.0.1:5000", "Address to host under")

var App *web.Mux
var Driver *structs.Driver


func init() {
    App = web.New()
    App.Use(middleware.Logger)
    App.Use(auth.Middleware)

    App.Handle("/d/*", func(c web.C, w http.ResponseWriter, r *http.Request) {
        target := fmt.Sprintf("http://%s",
            strings.SplitN(r.URL.Path, "/", 3)[2])

        req, err := http.NewRequest(r.Method, target, r.Body)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, err)
            return
        }

        resp, err := http.DefaultClient.Do(req)

        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, err)
            return
        }

        body, err := ioutil.ReadAll(resp.Body)

        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, err)
            return
        }

        w.WriteHeader(resp.StatusCode)
        w.Write(body)
    })

    App.Get("/vms", func(c web.C, w http.ResponseWriter, r *http.Request) {
        response, _ := json.Marshal(Driver.GetVMs())
        w.Write(response)
    })

    App.Get("/which", func(c web.C, w http.ResponseWriter, r *http.Request) {
        response, _ := json.Marshal(Driver.Name)
        w.Write(response)
    })

    App.Compile()
}


func main() {
    log.Printf(ansi.Color("Starting Beacon...", "white+b"))

    if !flag.Parsed() {
        flag.Parse()
    }

    Driver = drivers.Decide()

    log.Printf("Provider Interface: %s\n", ansi.Color(Driver.Name, "cyan+b"))
    log.Printf("Authentication Token: %s\n", ansi.Color(*auth.Token, "cyan+b"))

    graceful.HandleSignals()

    graceful.PreHook(func() {
        log.Printf(ansi.Color("Gracefully Shutting Down...", "white+b"))
    })
    graceful.PostHook(func() {
        log.Printf(ansi.Color("Done!", "white+b"))
    })

    defer graceful.Wait()

    http.Handle("/", App)
    log.Printf("Listening on %s", *address)

    if *pemFile != "" && *keyFile != "" {
        log.Printf("Setting up secure server...")
        err = graceful.ListenAndServeTLS(*address, *pemFile, *keyFile, http.DefaultServeMux)
    } else {
        log.Printf(ansi.Color("Setting up unsecure server...", "yellow+b"))
        err = graceful.ListenAndServe(*address, http.DefaultServeMux)
    }


    if err != nil {
        log.Fatal(err)
    }
}
