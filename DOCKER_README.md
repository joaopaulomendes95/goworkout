 # GoWorkout Application

This README provides instructions for running the GoWorkout application in both development and production environments.

## Development Environment

### Prerequisites

- Docker and Docker Compose installed
- Git (to clone the repository)

### Setup and Run

1. Clone the repository:

   ```bash
   git clone <repository-url>
   cd goworkout
     2. Create a .env file in the root directory with the following content:
   ```

 PORT=8080
APP_ENV=local
GOWORKOUT_DB_HOST=localhost
GOWORKOUT_DB_PORT=5432
GOWORKOUT_DB_DATABASE=goworkout
GOWORKOUT_DB_USERNAME=postgres
GOWORKOUT_DB_PASSWORD=postgres
GOWORKOUT_DB_SCHEMA=public
   3. Start the development environment:

 docker-compose -f docker-compose.dev.yml up
 Note: If you want to use the default filename, you can rename docker-compose.dev.yml to docker-compose.yml and then simply run:

 docker-compose up
 Docker Compose by default looks for a file named docker-compose.yml or docker-compose.yaml in the current directory.  4. Access the application:
▪ Frontend: <http://localhost:5173>  ▪ Backend API: <http://localhost:8080>  ▪ Database: PostgreSQL running on port 5432  5. To stop the development environment:

 docker-compose -f docker-compose.dev.yml down
 Or if using the default filename:

 docker-compose down

Development Notes
• The frontend uses SvelteKit with hot-reloading enabled  • The backend uses Go with Chi router  • Changes to the frontend and backend code will automatically reload
Production Environment
Prerequisites
• Docker and Docker Compose installed  • Git (to clone the repository)
Setup and Run 1. Clone the repository:

 git clone <repository-url>
cd goworkout
   2. Create a .env file in the root directory with the following content:

 PORT=8080
APP_ENV=production
GOWORKOUT_DB_HOST=localhost
GOWORKOUT_DB_PORT=5432
GOWORKOUT_DB_DATABASE=goworkout
GOWORKOUT_DB_USERNAME=postgres
GOWORKOUT_DB_PASSWORD=postgres
GOWORKOUT_DB_SCHEMA=public
   3. Start the production environment:

 docker-compose -f docker-compose.prod.yml up --build
   4. Access the application:
▪ Frontend: <http://localhost:5173>  ▪ Backend API: <http://localhost:8080>  ▪ Database: PostgreSQL running on port 5432  5. To stop the production environment:

 docker-compose -f docker-compose.prod.yml down
   6. To completely clean up (remove images, volumes, etc.):

 docker-compose -f docker-compose.prod.yml down --rmi all --volumes --remove-orphans

Production Notes
• The frontend is built for production and served using Node.js  • The backend is compiled to a binary for optimal performance  • The database data is persisted in a Docker volume
Differences Between Development and Production
• Development: Frontend and backend run in development mode with hot-reloading  • Production: Frontend is built and optimized, backend is compiled to a binary  • Both environments use the same ports and configuration for consistency
