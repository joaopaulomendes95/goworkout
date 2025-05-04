FROM golang:1.24.2-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/app/main.go

FROM alpine:3.20.1 AS prod
WORKDIR /app
COPY --from=build /app/main /app/main
EXPOSE ${PORT}
CMD ["./main"]

FROM node:20 AS frontend_builder
WORKDIR /frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/. .
RUN npm run build

# FROM node:23-slim AS frontend
# RUN npm install -g serve
# COPY --from=frontend_builder /frontend/.svelte-kit/output/client /app/dist
# EXPOSE 5173
# CMD ["serve", "-s", "/app", "-l", "5173"]

# Change just this stage to run dev instead of serve
FROM node:20-slim AS frontend
WORKDIR /app
# Copy everything instead of just the build output
COPY --from=frontend_builder /frontend /app
EXPOSE 5173
# Run dev server instead of serve static files
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]
