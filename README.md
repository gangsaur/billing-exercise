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

# Setup

Obviously for local dev only since it's an exercise repo.

1. Setup dependency using docker
```
cd deploy/dev
docker-compose up -d
```

2. Run the init script to create initial tables
```
go run cmd/script/run-sql/main.go 01-init.sql
```

3. Run the API service
```
go run cmd/api/main.go
```
