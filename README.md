# TiDB + Gin + GORM CRUD Demo

This is a simple Create-Read-Update-Delete (CRUD) example project using TiDB as the database, the Gin web framework, and GORM as the ORM.

## Project Structure

```
tidb-gin-demo/
├── config/
│   └── database.go      # Database configuration and connection
├── controllers/
│   └── user_controller.go  # User controller (CRUD handlers)
├── models/
│   └── user.go          # User model and request structs
├── main.go              # Application entry point
├── go.mod               # Go module file
└── README.md            # Project description (Chinese)
```

## Requirements

- Go 1.21+
- TiDB database (local or remote)

## Install dependencies

```bash
go mod tidy
```

## Configure TiDB connection

Modify the connection string in `config/database.go`:

```go
dsn := "username:password@tcp(host:port)/database?charset=utf8mb4&parseTime=True&loc=Local"
```

### TiDB Cloud TLS (optional)

If you use TiDB Cloud (or another gateway that requires TLS), see the configuration examples below. The project supports enabling TLS via environment variables in `config/database.go`, and optionally loading a custom CA.

- **Environment variables** (replace with your real values):
  - `DB_USER`: database username
  - `DB_PASS`: database password
  - `DB_HOST`: database host (for example `gateway01.us-west-2.prod.aws.tidbcloud.com`)
  - `DB_PORT`: port (TiDB Cloud often uses `4000`)
  - `DB_NAME`: database name
  - `TIDB_TLS`: set to `true` to enable TLS (disabled by default)
  - `TIDB_TLS_SERVERNAME`: optional, TLS server name for certificate verification (the gateway host)
  - `TIDB_TLS_CA`: optional, path to a CA PEM file; leave empty to use system root CAs

- **Single-ENV DSN example (macOS/zsh)**:

```bash
# Export a full DSN once (includes &tls=tidb)
export DB_DSN='user.root:<PASSWORD>@tcp(host:port)/test?charset=utf8mb4&parseTime=True&loc=Local&tls=tidb'

go run main.go
```

Alternatively you can export components separately:

```bash
export DB_USER="xxxxx.root"
export DB_PASS="abcdefgxxxxxxmn"
export DB_HOST="gateway99.uk-east-8.prod.aws.tidbcloud.com"
export DB_PORT="4000"
export DB_NAME="test"
# Optional if DSN already contains tls=tidb
export TIDB_TLS="true"
export TIDB_TLS_SERVERNAME="gateway99.uk-east-8.prod.aws.tidbcloud.com"
# If you use a self-signed CA:
# export TIDB_TLS_CA="/path/to/ca.pem"

go run main.go
```

- **Equivalent DSN (manual form)**:

```
user:pass@tcp(host:port)/test?charset=utf8mb4&parseTime=True&loc=Local&tls=tidb
```

Note: Do not commit real credentials to source control or logs. If a password contains special characters, either escape them correctly or pass the full DSN via `DB_DSN`.

## Suppress common Gin warnings

You may see warnings at runtime such as:

- `Creating an Engine instance with the Logger and Recovery middleware already attached.`
- `Running in "debug" mode. Switch to "release" mode in production.`
- `You trusted all proxies, this is NOT safe. We recommend you to set a value.`

Remedies:

- Switch to production mode (macOS/zsh):

```bash
export GIN_MODE=release
# or in code: gin.SetMode(gin.ReleaseMode)
```

- The code uses `gin.New()` and manually applies `gin.Logger()` and `gin.Recovery()` to avoid double-attaching middleware.

- Configure trusted proxies (the default example trusts `127.0.0.1`). If your service sits behind proxies/load balancers, set a comma-separated list of IPs/CIDRs:

```bash
export TRUSTED_PROXIES="10.0.0.0/8,192.168.1.1"
```

If you are unsure, keep the default during development and adjust for production topology.

## Run the project

```bash
go run main.go
```

The service runs at `http://localhost:8080` by default.

## API Endpoints

### Health Check
- **GET** `/health` — check service status

### User Management
- **POST** `/api/users/` — create a user
- **GET** `/api/users/` — get all users
- **GET** `/api/users/:id` — get a single user
- **PUT** `/api/users/:id` — update a user
- **DELETE** `/api/users/:id` — delete a user

## Examples

### 1. Create user
```bash
curl -X POST http://localhost:8080/api/users/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "johndoe@example.com",
    "age": 25
  }'
```

### 2. Get all users
```bash
curl http://localhost:8080/api/users/
```

### 3. Get a single user
```bash
curl http://localhost:8080/api/users/1
```

### 4. Update a user
```bash
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "age": 26
  }'
```

### 5. Delete a user
```bash
curl -X DELETE http://localhost:8080/api/users/1
```

## Data Model

### User model
- `id`: primary key, auto-increment
- `name`: user name (required)
- `email`: email (required, unique)
- `age`: age
- `created_at`: creation time
- `updated_at`: update time

## Notes

1. Make sure TiDB is running
2. Ensure database connection settings are correct
3. The first run will auto-create database tables
4. All API responses are JSON

## Error handling

- 400 Bad Request: invalid request parameters
- 404 Not Found: resource not found
- 500 Internal Server Error: server error
