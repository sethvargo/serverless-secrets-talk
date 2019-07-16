// Copyright 2019 Seth Vargo
// Copyright 2019 Google, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/gomodule/redigo/redis"
	_ "github.com/sethvargo/go-malice"
)

var (
	appEnv    = getEnvOrDefault("ENV", "development")
	appPort   = getEnvOrDefault("PORT", "8080")
	redisHost = getEnvOrDefault("REDIS_HOST", "127.0.0.1")
	redisPort = getEnvOrDefault("REDIS_PORT", "6379")
	redisPass = getEnvOrDefault("REDIS_PASS", "")
)

type app struct {
	server    *http.Server
	redisPool *redis.Pool
	env       string
}

func main() {
	redisPool := &redis.Pool{
		MaxIdle:     3,
		MaxActive:   10,
		IdleTimeout: 30 * time.Second,
		Dial: func() (redis.Conn, error) {
			addr := redisHost + ":" + redisPort
			return redis.Dial("tcp", addr, redis.DialPassword(redisPass))
		},
	}

	a := &app{
		server:    &http.Server{},
		redisPool: redisPool,
		env:       appEnv,
	}

	http.HandleFunc("/favicon.ico", a.notFoundHandler())
	http.HandleFunc("/reset-counter", a.resetCounterHandler())
	http.HandleFunc("/", a.indexHandler())
	log.Fatal(http.ListenAndServe(":"+appPort, nil))
}

func (a *app) indexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn := a.redisPool.Get()
		defer conn.Close()

		count, err := redis.Int(conn.Do("INCR", "visits"))
		if err != nil {
			a.handleError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")

		result := processTemplate(htmlIndex, struct{ Count int }{count})
		fmt.Fprintf(w, result)
	}
}

func (a *app) resetCounterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn := a.redisPool.Get()
		defer conn.Close()

		val := r.URL.Query().Get("count")
		if val == "" {
			val = "0"
		}

		_, err := conn.Do("SET", "visits", val)
		if err != nil {
			a.handleError(w, r, err)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (a *app) handleError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("[ERR] %s", err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/html")

	if a.env == "development" {
		env := os.Environ()
		sort.Strings(env)

		result := processTemplate(htmlDevErr, struct {
			Error         error
			Request       *http.Request
			XForwardedFor string
			Env           []string
		}{
			Error:         err,
			Request:       r,
			XForwardedFor: r.Header.Get("x-forwarded-for"),
			Env:           env,
		})
		fmt.Fprintf(w, result)
	} else {
		fmt.Fprint(w, htmlProdErr)
	}
}

func (a *app) notFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}
}

func getEnvOrDefault(key, def string) string {
	v := os.Getenv(key)
	if v != "" {
		return v
	}
	return def
}
