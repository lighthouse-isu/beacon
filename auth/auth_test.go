package auth

import (
    "fmt"
    "testing"
    "net/http"
    "net/http/httptest"
)


func TestGenerateToken(t *testing.T) {
    testToken := GenerateToken(42)
    if len(testToken) != 84 {
        t.Fail()
    }
}

func okHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "OK")
}

func TestMiddlewareUnauthed(t *testing.T) {
    *Token = "abc123"

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
    Middleware(nil, http.HandlerFunc(okHandler)).ServeHTTP(w, req)

    if w.Body.String() == "OK" {
        t.Fail()
    }
}

func TestMiddlewareIncorrectAuthed(t *testing.T) {
    *Token = "abc123"

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
    req.Header["Token"] = []string{"this is wrong"}
    Middleware(nil, http.HandlerFunc(okHandler)).ServeHTTP(w, req)

    if w.Body.String() == "OK" {
        t.Fail()
    }
}

func TestMiddlewareAuthed(t *testing.T) {
    *Token = "abc123"

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
    req.Header["Token"] = []string{"abc123"}
    Middleware(nil, http.HandlerFunc(okHandler)).ServeHTTP(w, req)

    if w.Body.String() != "OK" {
        t.Fail()
    }
}