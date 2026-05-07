package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"gangsaur.com/billing-exercise/internal/middleware"
)

type ApiServer struct {
	server *http.Server
}

func NewApiServer() *ApiServer {
	// Init mux and server
	mux := createMux()

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "9000" // Simply default to 9000
	}

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("Received %s, shutting down", sig)
			err := srv.Shutdown(context.Background())
			if err != nil {
				log.Printf("Error when shutting down server: %s", err.Error())
			}
		}
	}()

	return &ApiServer{
		server: &srv,
	}
}

func (s ApiServer) GetServerAddress() string {
	return s.server.Addr
}

func (s ApiServer) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func createMux() *http.ServeMux {
	mux := http.NewServeMux()
	basicChain := middleware.MiddlewareChain{middleware.RequestMiddleware, middleware.LogMiddleware}

	mux.HandleFunc("GET /h", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Temporairly for testing middleware
	mux.Handle("GET /loan/{id}", basicChain.ThenFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}))

	return mux
}
