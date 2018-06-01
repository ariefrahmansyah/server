package main

import (
	"context"
	"net/http"
	"time"

	"github.com/ariefrahmansyah/server"
)

func main() {
	ctxWeb := context.Background()

	serverOptions := &server.Options{
		ListenAddress:  "localhost:8080",
		MaxConnections: 512,
		ReadTimeout:    10 * time.Second,
	}

	webServer := server.New(nil, serverOptions)
	webServer.Router().Get("/simple", simple)

	if err := webServer.Run(ctxWeb); err != nil {
		panic(err)
	}
}

func simple(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("simple"))
}
