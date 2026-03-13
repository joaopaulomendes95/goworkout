# Go Svelte Fullstack

A minimal fullstack template with Go backend and SvelteKit frontend.

## Quick Start

```bash
# Start development
docker compose --profile dev up --build
```

This starts:
- **PostgreSQL** on port 5432
- **Go backend** on port 8080
- **SvelteKit frontend** on port 5173

## Project Structure

```
/backend          # Go API
/frontend         # SvelteKit app
docker-compose.yml
Makefile
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/users` | Register |
| POST | `/api/tokens/authentication` | Login |
| GET | `/api/me` | Get current user |
| GET | `/api/hello` | Protected example |
| GET | `/health` | Health check |
