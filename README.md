# Azure Client Go

A clean-architecture Go service with:
- Gorilla Mux HTTP API
- GORM (MySQL) with OpenTelemetry tracing
- Dependency Injection (DI) for repositories and services
- Custom error handling and middleware
- Modular structure for easy extension
- Example Azure and HTTP client integrations
- CI/CD via GitHub Actions

## Project Structure
```
azure-client/
├── main.go
├── client/
│   ├── azure/
│   └── http/
├── config/
├── internal/
│   ├── controller/
│   ├── model/
│   ├── repository/
│   ├── service/
│   ├── middleware/
│   ├── errs/
│   └── otel/
├── .github/
│   └── workflows/
│       └── ci.yml
├── go.mod
├── go.sum
└── README.md
```

## Features
- **CRUD API** for `Member` entity
- **OpenTelemetry** tracing for HTTP, service, and DB layers
- **GORM** for MySQL ORM
- **Azure/HTTP client** integration examples
- **Centralized error handling**
- **CI/CD**: Lint, test, build on every push/PR

## Quick Start
1. **Clone the repo:**
   ```sh
   git clone <your-repo-url>
   cd azure-client
   ```
2. **Configure environment:**
   - Copy `.env.example` to `.env` and fill in DB/Azure config
3. **Run locally:**
   ```sh
   go run main.go
   ```
4. **Test API:**
   - Use Postman/curl to hit `/members` endpoints
5. **View traces:**
   - Run Jaeger locally (see `otel-test-instructions.md`)
   - Open [http://localhost:16686](http://localhost:16686)

## CI/CD
- See `.github/workflows/ci.yml` for pipeline details
- Lint, staticcheck, test, and build run on every push/PR

## Extending
- Add new domains: create new files in `internal/model`, `repository`, `service`, `controller`
- Add new clients: add to `client/`
- Add new middleware: add to `internal/middleware/`

## License
MIT
