# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Better WAPI is a RESTful API wrapper around the WEDOS DNS API (WAPI). It provides standardized CRUD operations for DNS records that the original WAPI lacks, using JWT authentication and modern API patterns.

## Development Commands

### Running the application
```bash
go run .
```

The API runs on port 8000 by default. Swagger UI is accessible at http://localhost:8000/docs/index.html (root path redirects there automatically).

### Building
```bash
go build -o app ./main.go
```

### Swagger documentation regeneration
```bash
swag init
```

Run this after modifying API handlers with godoc annotations in src/api/handlers/.

### Docker
```bash
docker build -t better-wapi:latest .
docker run -d -p 8083:8000 --env-file .env better-wapi:latest
```

### Testing
This project currently has no automated tests.

### Upgrading Go version
When upgrading the Go version, update it in the following files:
- `go.mod` (line 3: `go 1.xx`)
- `Dockerfile` (line 1: `FROM golang:1.xx-alpine`)
- `.github/workflows/ci.yaml` (line 26: `go-version: "1.xx"`)

After updating, run:
```bash
go get -u           # Update dependencies to latest minor versions
go mod tidy         # Clean up go.mod and go.sum
go build -o app ./main.go  # Verify build succeeds
```

## Architecture

### Core flow
1. **Authentication**: JWT-based auth using BW_USER credentials (src/services/auth.go)
   - Obtain token via POST /token or POST /api/auth/token with login/secret credentials
   - Token must be passed as Bearer token in Authorization header for all protected endpoints
2. **Request handling**: Gin handlers in src/api/handlers/ receive RESTful requests
   - All /v1/domain and /v2/domain routes require JWT authentication
3. **WAPI integration**: IntegrationService (src/services/integration.go) translates to WAPI commands
4. **WAPI authentication**: Time-based SHA1 token generated using WAPI credentials + Prague timezone hour
   - Token format: SHA1(username + SHA1(password) + current_hour_in_prague)
   - Generated fresh for each WAPI request via getApiToken() function

### Key components

**IntegrationService** (src/services/integration.go)
- Core service that communicates with WEDOS API
- Handles WAPI authentication via time-based SHA1 token (Prague timezone)
- All DNS operations go through `makeRequest()` method
- Supports optional autocommit parameter for immediate DNS changes

**API versioning**
- v1 endpoints: Use subdomain strings for lookups (e.g., `/v1/domain/{domain}/record`)
  - Create: POST /v1/domain/{domain}/record
  - Update: PUT /v1/domain/{domain}/record (finds record by subdomain)
  - Delete: DELETE /v1/domain/{domain}/record (finds record by subdomain)
  - Get all: GET /v1/domain/{domain}/info
  - Get one: GET /v1/domain/{domain}/{subdomain}/info
  - Commit: POST /v1/domain/{domain}/commit
- v2 endpoints: Use numeric record IDs from path (e.g., `/v2/domain/{domain}/record/{id}`)
  - Update: PUT /v2/domain/{domain}/record/{id}
  - Delete: DELETE /v2/domain/{domain}/record/{id}
  - v2 returns 204 No Content on success (v1 returns 200/201)

**Request models** (src/api/models/requests.go)
- SaveRowRequest (v1): Requires subdomain string in request body
- SaveRowRequestV2 (v2): Uses record ID from URL path parameter
- Both support fields: TTL, Type, Data, Autocommit
- Default values (applied via creasty/defaults package):
  - TTL: 3600
  - Type: "A"
  - Autocommit: false
- DeleteRowRequest/DeleteRowRequestV2: Support Autocommit field (default: false)

**Middleware**
- Authorization middleware in src/api/middleware/authorize.go validates JWT tokens
- Applied to all /v1 and /v2 domain routes via router groups in src/api/routes.go
- CORS is configured to allow all origins with credentials in main.go

### Environment configuration

Required environment variables (see .env.example):
- `BW_USER_LOGIN` / `BW_USER_SECRET`: API authentication credentials
- `BW_WAPI_USERNAME` / `BW_WAPI_PASSWORD`: WEDOS WAPI credentials
- `BW_JSON_WEB_KEY`: JWT signing key
- `BW_BASE_URL`: API base URL for Swagger docs
- `BW_USE_LOGFILE`: Enable file logging

**IMPORTANT**: WEDOS requires whitelisting the host IP address in their management dashboard.

## Certbot integration

The `tools/certbot/certbot_renew_hook.py` script provides DNS-01 challenge support for Let's Encrypt certificate renewal via the Better WAPI API. It creates and cleans up TXT records for ACME challenges.
