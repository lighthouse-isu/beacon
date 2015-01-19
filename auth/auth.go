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

package auth

import (
    "flag"
    "net/http"
    "encoding/hex"
    "crypto/rand"

    "github.com/zenazn/goji/web"
)

var Token = flag.String("token", GenerateToken(32), "Predefined auth token")

func Middleware(c *web.C, h http.Handler) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        requestToken, ok := r.Header["Token"]
        if ok && requestToken[0] == *Token {
            h.ServeHTTP(w, r)
        } else {
            w.WriteHeader(401)
        }
    }

    return http.HandlerFunc(fn)
}

func GenerateToken(size int) string {
    token := make([]byte, size)
    rand.Read(token)
    return hex.EncodeToString(token)
}
