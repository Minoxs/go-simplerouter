package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/minoxs/simplerouter"
)

func middlewareExample(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now().Second()
	log.Println("Request to:", r.URL.Path)
	next(w, r)
	log.Println("Request finished:", r.URL.Path, "Delay (s):", time.Now().Second()-start)
}

var pingCount = 0

func pong(w http.ResponseWriter, r *http.Request) {
	pingCount += 1
	_, _ = fmt.Fprint(w, "PONG")
}

func count(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, pingCount)
}

func reset(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "Count reset", pingCount)
	pingCount = 0
}

func randomTime(w http.ResponseWriter, r *http.Request) {
	n := rand.Intn(10)
	t, _ := time.ParseDuration(fmt.Sprintf("%ds", n))
	time.Sleep(t)
	_, _ = fmt.Fprint(w, "Hello, World!")
}

func main() {
	r := simplerouter.New().
		Get("/ping", pong).
		Post("/ping", pong)

	r.Prefix("/status").
		Get("/count", count).
		Delete("/reset", reset)

	r.Use(middlewareExample).
		Get("/test", randomTime)

	server := &http.Server{
		Addr:    ":80",
		Handler: r,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
