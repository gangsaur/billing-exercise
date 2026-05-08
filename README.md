# billing-exercise

Simple billing system exercise.

# Commit Tag

Tags for grouping commit content, not a strict requirements.

- [Documentation]: Changes related to documentation only
- [Feature]: Changes related to functionality
- [Refactor]: Changes related to restructuring which does not or mostly not affect functionality
- [Build]: Changes related to build, setup, and deployment
- [Test]: Changes related to test only

# Requirements

- Go
- Docker
- Mockery (for mock generation only, not needed for running)

# Setup

Obviously for local dev only since it's an exercise repo.

1. Setup dependency using docker
```
cd deploy/dev
docker-compose up -d
```

2. Run the init script from the project root directory to create initial tables. Also seed it if needed.
```
go run cmd/script/run-sql/main.go 01-init.sql

# Only if you want mock starting data
go run cmd/script/run-sql/main.go 02-seed.sql
```

3. Run the API service
```
go run cmd/api/main.go
```

# Notes

- GetOutstanding: Implemented as part of `GetLoan` in `internal/service/loan.go` (OutstandingAmount is saved in DB)
- IsDelinquent: Implemented as part of `GetUser` in `internal/service/user.go` (Calculated by fetching user's loan and loanPayments)
- MakePayment: Implemented as `PayLoan` in `internal/service/loan.go`
