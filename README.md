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

2. Copy the env.sample to .env and fill it
```
# Go back to root, assuming 1st step is run
cd ../../

cp env.sample .env
```

3. Run the init script from the project root directory to create initial tables. Also seed it if needed.
```
# Run the init script and mock starting data
go run cmd/script/run-sql/main.go 01-init.sql
go run cmd/script/run-sql/main.go 02-seed.sql
```

4. Run the API service
```
go run cmd/api/main.go
```

# Notes

- `GetOutstanding`: Implemented as part of `GetLoan` in `internal/service/loan.go` (OutstandingAmount is saved in DB)
- `IsDelinquent`: Implemented as part of `GetUser` in `internal/service/user.go` (Calculated by fetching user's loan and loanPayments)
- `MakePayment`: Implemented as `PayLoan` in `internal/service/loan.go`

# Request Sample

- `GetOutstanding` via `GET /loan/{id}`
```
curl -X GET http://localhost:9990/loan/1
curl -X GET http://localhost:9990/loan/2
```
= `IsDelinquent` via `POST /loan/{id}/pay`
```
curl -X POST -i --data '{"amount": 110000}' http://localhost:9990/loan/1/pay
```
- `MakePayment` via `GET /user/{id}`
```
curl -X GET http://localhost:9990/user/1
curl -X GET http://localhost:9990/user/2
```