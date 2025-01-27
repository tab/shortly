# Shortly

## Overview
Shortly is a Go application designed to shorten long URLs

## Setup Instructions

### Prerequisites

**Docker** and **Docker Compose** installed on your machine

### Build

```sh
docker-compose build
```

### Database setup

```sh
docker-compose up -d database
```

```sh
GO_ENV=development make db:create db:migrate
```

#### Database management commands

The Makefile provides commands for managing the database:

    Create Database: make db:create
    Drop Database: make db:drop
    Apply Migrations: make db:migrate
    Check Migration Status: make db:migrate:status
    Rollback Last Migration: make db:rollback
    Dump Schema: make db:schema:dump
    Load Schema: make db:schema:load

### Run application

```sh
docker-compose up backend
```

### API Documentation

Check api/swagger.yml for the API documentation

### Development

Check `.env.development` for the environment variables
and `.env.test` for the test environment variables

#### Prepare test environment

```sh
GO_ENV=test make db:create db:schema:load
```

Unit Tests:

```sh
GO_ENV=test make test
```

Test Coverage:

```sh
make coverage
```

Code Vetting:

```sh
make vet
```

Linting:

```sh
make lint
```

### Profiling

Generate payload with wrk tool:

```sh
make benchmark:payload
```

Create CPU and Memory profiles:

```sh
make pprof:cpu
make pprof:mem
```

Compare profiles:

```sh
make pprof:cpu:diff
make pprof:mem:diff
```
