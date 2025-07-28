# Test_effective_mobile

## Quick Start

1. Copy the environment file:

   ```sh
   cp .env.example .env
   ```

2. Start the Docker containers:

```
 docker compose up --build -d

```

3. Access the application:
   Swagger UI: http://127.0.0.1:4000/swagger/index.html#/

## Migration Commands

| Command                      | Description                         |
| ---------------------------- | ----------------------------------- |
| `go run cmd/migrate.go up`   | Apply all pending migrations        |
| `go run cmd/migrate.go down` | Roll back the most recent migration |
