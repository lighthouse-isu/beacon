package structs

import (
    "fmt"
    "strings"
    "testing"
    "net"
    "net/http"
    "net/http/httptest"
)


func TestPingDocker(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "OK")
    }))
    defer ts.Close()

    host, port, _ := net.SplitHostPort(
        strings.Replace(ts.URL, "http://", "", 1))

    vm := &VM{
        Name: "test",
        Address: host,
        Port: port,
        Version: "v1.23",
        CanAccessDocker: false,
    }
    result := vm.PingDocker()

    if !result {
        t.Fail()
    }
}

func TestPingDockerDNE(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(404)
    }))
    defer ts.Close()

    host, port, _ := net.SplitHostPort(
        strings.Replace(ts.URL, "http://", "", 1))

    vm := &VM{
        Name: "test",
        Address: host,
        Port: port,
        Version: "v1.23",
        CanAccessDocker: false,
    }
    result := vm.PingDocker()

    if result {
        t.Fail()
    }
}

func TestGetDockerVersion(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "{\"ApiVersion\": \"1.14\"}")
    }))
    defer ts.Close()

    host, port, _ := net.SplitHostPort(
        strings.Replace(ts.URL, "http://", "", 1))

    vm := &VM{
        Name: "test",
        Address: host,
        Port: port,
        Version: "v1.23",
        CanAccessDocker: false,
    }

    result, err := vm.GetDockerVersion()

    if result != "v1.14" || err != nil {
        t.Fail()
    }
}
