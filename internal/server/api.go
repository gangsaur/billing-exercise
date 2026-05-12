package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"gangsaur.com/billing-exercise/internal/handler"
	"gangsaur.com/billing-exercise/internal/middleware"
	"gangsaur.com/billing-exercise/internal/repository/db/psql"
	"gangsaur.com/billing-exercise/internal/service"
)

type ApiServer struct {
	server *http.Server
}

func NewApiServer() *ApiServer {
	// Init dependencies
	dsn := os.Getenv("PSQL_DSN")
	if dsn == "" {
		log.Fatal("PSQL_DSN is required")
	}
	psql, err := psql.NewPsql(dsn)
	if err != nil {
		log.Fatalf("Failed creating PSQL pool: %v", err.Error())
	}

	// Init mux and server
	mux := createMux(psql)

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

			psql.CloseConnection()

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

func createMux(p *psql.Psql) *http.ServeMux {
	mux := http.NewServeMux()
	basicChain := middleware.MiddlewareChain{middleware.RequestMiddleware, middleware.LogMiddleware}

	loanService := service.NewLoanService(p)
	loanHandler := handler.NewLoanHandler(loanService)

	userSevice := service.NewUserService(p)
	userHandler := handler.NewUserHandler(userSevice)

	mux.HandleFunc("GET /h", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Loan
	mux.Handle("GET /loan/{id}", basicChain.ThenFunc(loanHandler.GetLoan()))
	mux.Handle("GET /loan/{id}/payments", basicChain.ThenFunc(loanHandler.GetLoanAndLoanPayments()))
	mux.Handle("POST /loan/{id}/pay", basicChain.ThenFunc(loanHandler.PayLoan()))

	// User
	mux.Handle("GET /user/{id}", basicChain.ThenFunc(userHandler.GetUser()))

	return mux
}
