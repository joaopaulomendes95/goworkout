# GoWorkout Application

Instructions for running GoWorkout in development and production.

---

## Development

### Prerequisites

- Docker & Docker Compose
- Git

### Setup

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd goworkout
   ```

2. **Create `.env` in the root:**

   ```
   PORT=8080
   APP_ENV=local
   GOWORKOUT_DB_HOST=localhost
   GOWORKOUT_DB_PORT=5432
   GOWORKOUT_DB_DATABASE=goworkout
   GOWORKOUT_DB_USERNAME=postgres
   GOWORKOUT_DB_PASSWORD=postgres
   GOWORKOUT_DB_SCHEMA=public
   ```

3. **Frontend API calls:**  
   Use Docker service name (`app`) instead of `localhost`:

   ```ts
   // Use for Docker development
   const API = "http://app:8080";
   // Not: const API = 'http://localhost:8080';
   ```

4. **Start development:**

   - From root:

     ```bash
     docker-compose -f docker-compose.dev.yml up
     ```

   - From frontend:

     ```bash
     npm run start
     ```

   - If using `docker-compose.yml`:

     ```bash
     docker-compose up
     ```

5. **Access:**

   - Frontend: [http://localhost:5173](http://localhost:5173)
   - Backend: [http://localhost:8080](http://localhost:8080)
   - DB: PostgreSQL on port 5432

6. **Stop:**

   ```bash
   docker-compose -f docker-compose.dev.yml down
   # or
   docker-compose down
   ```

#### Notes

- Frontend: SvelteKit (hot-reloading)
- Backend: Go (Chi router)
- Code changes auto-reload
- Use `http://app:8080` for API calls from frontend (in Docker)
- Use `http://localhost:8080` to access backend directly

---

## Production

### Prerequisites

- Docker & Docker Compose
- Git

### Setup

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd goworkout
   ```

2. **Create `.env` in the root:**

   ```
   PORT=8080
   APP_ENV=production
   GOWORKOUT_DB_HOST=localhost
   GOWORKOUT_DB_PORT=5432
   GOWORKOUT_DB_DATABASE=goworkout
   GOWORKOUT_DB_USERNAME=postgres
   GOWORKOUT_DB_PASSWORD=postgres
   GOWORKOUT_DB_SCHEMA=public
   ```

3. **Start production:**

   - From root:

     ```bash
     docker-compose -f docker-compose.prod.yml up --build
     ```

   - From frontend:

     ```bash
     npm run start:prod
     ```

4. **Access:**

   - Frontend: [http://localhost:5173](http://localhost:5173)
   - Backend: [http://localhost:8080](http://localhost:8080)
   - DB: PostgreSQL on port 5432

5. **Stop:**

   ```bash
   docker-compose -f docker-compose.prod.yml down
   ```

6. **Clean up everything:**

   ```bash
   docker-compose -f docker-compose.prod.yml down --rmi all --volumes --remove-orphans
   ```

#### Notes

- Frontend: Built for production, served via Node.js
- Backend: Compiled Go binary
- Database persists in Docker volume
- API routing handled automatically

---

## Dev vs Production

| Aspect             | Development                            | Production                        |
| ------------------ | -------------------------------------- | --------------------------------- |
| Frontend           | SvelteKit, hot-reloading               | Built & optimized, served by Node |
| Backend            | Go (dev mode, Chi router)              | Compiled Go binary                |
| API URL (frontend) | <http://app:8080> (Docker service)       | Handled automatically             |
| DB Persistence     | Docker volume                          | Docker volume                     |
| Ports              | 5173 (frontend), 8080 (API), 5432 (DB) | Same as development               |

---
