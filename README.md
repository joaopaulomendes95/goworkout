# goworkout

This project is a basic web application to track workouts.
Its a study subject to get familiar with GO, and frontend in general.

Special thanks to:

- <https://github.com/Melkeydev/go-blueprint>
- <https://github.com/go-chi/chi>
- <https://github.com/air-verse/air>
- <https://github.com/sveltejs/kit>
- <https://github.com/tailwindlabs/tailwindcss>
- <https://github.com/neovim/neovim>

Some way or another, these projects were either used as a reference or inspiration

## Technology Stack

- **Backend**: Go with Chi router
- **Frontend**: Svelte + Tailwind CSS
- **Database**: PostgreSQL with Goose migrations
- **Authentication**: Custom token-based auth

## Project Structure

- `/cmd` - Application entry point
  - `/app` - Application main package
- `/frontend` - Web client code
  - `/src` - Frontend source code
  - `/static` - Static assets
- `/internal` - Private application code
  - `/api` - REST API handlers
  - `/database` - Database connection management
  - `/server` - HTTP server configuration
  - `/store` - Data access layer
  - `/utils` - Helper functions
- `/migrations` - PostgreSQL schema migrations

## Requirements

- Go 1.24 or higher
- PostgreSQL 14+
- Docker (optional)

## Features

- User authentication and management
- Create and track workouts
- Log exercises with sets, reps and weights
- Responsive web interface

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.
See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests

```bash
make all
```

Build the application

```bash
make build
```

Run the application

```bash
make run
```

Create DB container

```bash
make docker-run
```

Shutdown DB Container

```bash
make docker-down
```

DB Integrations Test:

```bash
make itest
```

Live reload the application:

```bash
make watch
```

Run the test suite:

```bash
make test
```

Clean up binary from the last build:

```bash
make clean
```
