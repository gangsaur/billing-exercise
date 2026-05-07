package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatalf("Missing arguments, go run example: `go run cmd/script/run-sql <filename>`")
	}
	filename := args[1]

	_ = godotenv.Load()
	dsn := os.Getenv("PSQL_URL")
	if dsn == "" {
		log.Fatal("PSQL_URL is required")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed creating pool: %v", err.Error())
	}
	defer pool.Close()

	body, err := os.ReadFile(fmt.Sprintf("cmd/script/run-sql/%v", filename))
	if err != nil {
		log.Fatalf("Failed reading file: %v", err.Error())
	}

	_, err = pool.Exec(ctx, string(body))
	if err != nil {
		log.Fatalf("Failed running the file cmd/script/run-sql/%v: %v", filename, err.Error())
	}

	log.Printf("Successfully run the file: cmd/script/run-sql/%v\n", filename)
}
