## Getting Started

```bash
# Copy .env.template to .env
cp .env.template .env

# Start docker containers
docker compose up -d

# migrate
cd cmd/migrations && goose postgres "host=localhost port=5432 user=gl password=gl dbname=gl sslmode=disable" up migrations

# seed
go run ./cmd/seed