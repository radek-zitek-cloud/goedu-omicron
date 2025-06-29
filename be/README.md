# GoEdu Control Testing Platform - Backend

This is the backend API server for the GoEdu Control Testing Platform, built with Go following clean architecture principles. The platform provides compliance control testing workflows for financial institutions with comprehensive audit trails and evidence management.

## 🏗️ Architecture Overview

The backend follows clean architecture patterns with clear separation of concerns:

```
cmd/                    # Application entry points
├── server/            # Main HTTP server
├── migrate/           # Database migration tool
└── seed/              # Database seeding tool

internal/              # Private application code
├── config/           # Configuration management
├── handlers/         # HTTP handlers (presentation layer)
│   ├── graphql/     # GraphQL resolvers
│   ├── rest/        # REST API handlers
│   └── ws/          # WebSocket handlers
├── services/        # Business logic layer
│   ├── control/     # Control management
│   ├── testing/     # Testing workflow
│   ├── evidence/    # Evidence collection
│   ├── user/        # User management
│   ├── audit/       # Audit logging
│   ├── notification/# Notifications
│   └── storage/     # File storage
├── repositories/    # Data access layer
│   ├── mongo/       # MongoDB repositories
│   └── redis/       # Redis repositories
├── models/          # Domain models
├── middleware/      # HTTP middleware
└── utils/           # Utility functions

pkg/                  # Public/reusable packages
├── logger/          # Structured logging
├── database/        # Database connection
├── cache/           # Redis cache
└── storage/         # Object storage

scripts/             # Development and deployment scripts
deployments/         # Deployment configurations
```

## 🚀 Features

### Core Functionality
- **Clean Architecture**: Separation of concerns with clear dependency boundaries
- **Environment Configuration**: Flexible configuration via environment variables and files
- **Structured Logging**: Comprehensive audit trails with correlation IDs
- **Health Checks**: Application and dependency health monitoring
- **Graceful Shutdown**: Proper resource cleanup during application termination

### Infrastructure
- **MongoDB Integration**: Connection pooling and automated indexing
- **Redis Caching**: Performance optimization with caching layer
- **Object Storage**: File management with MinIO integration
- **CORS Support**: Cross-origin resource sharing for frontend integration

### Security & Compliance
- **Audit Logging**: Detailed audit trails for regulatory compliance
- **Request Correlation**: Distributed tracing capabilities
- **Input Validation**: Comprehensive data validation
- **Security Headers**: Proper HTTP security headers

## 🛠️ Development Setup

### Prerequisites

- **Go 1.21+**: [Install Go](https://golang.org/doc/install)
- **MongoDB**: For primary data storage
- **Redis**: For caching (optional for development)
- **MinIO**: For object storage (optional for development)

### Quick Start

1. **Clone and Setup**:
   ```bash
   git clone <repository-url>
   cd goedu-omicron/be
   ./scripts/setup-dev.sh
   ```

2. **Configure Environment**:
   ```bash
   cp .env.template .env
   # Edit .env with your configuration
   ```

3. **Run the Application**:
   ```bash
   # Development mode with hot reload
   go run ./cmd/server
   
   # Or build and run
   go build -o bin/server ./cmd/server
   ./bin/server
   ```

4. **Verify Setup**:
   ```bash
   curl http://localhost:8080/health
   ```

### Configuration

The application uses a hierarchical configuration system:

1. **Default Values**: Sensible defaults for development
2. **Configuration File**: `config.yaml` (optional)
3. **Environment Variables**: Override any setting (highest priority)

#### Key Configuration Sections

```yaml
app:
  name: "GoEdu Control Testing Platform"
  port: 8080
  environment: "development"

database:
  uri: "mongodb://localhost:27017"
  database: "goedu"
  max_pool_size: 100

cache:
  host: "localhost"
  port: 6379

auth:
  jwt_secret: "your-secret-key"
  jwt_expiration: "24h"

logger:
  level: "info"
  environment: "development"
```

#### Environment Variables

All configuration can be overridden with environment variables using the `GOEDU_` prefix:

```bash
GOEDU_APP_PORT=9000
GOEDU_DATABASE_URI="mongodb://prod-db:27017"
GOEDU_LOGGER_LEVEL="error"
```

## 🗄️ Database

### MongoDB Collections

The application creates the following collections with optimized indexes:

- **organizations**: Client organizations with multi-tenancy
- **users**: User accounts with role-based access
- **controls**: Compliance controls and testing procedures
- **testing_cycles**: Testing periods and progress tracking
- **evidence_requests**: Evidence collection workflow
- **audit_logs**: Comprehensive audit trail

### Automatic Indexing

The application automatically creates optimized indexes on startup:

```go
// Example: Control collection indexes
{Keys: {"organization_id": 1, "control_id": 1}, Unique: true}
{Keys: {"framework": 1, "category": 1}}
{Keys: {"created_at": -1}}
```

## 📊 Monitoring & Observability

### Health Endpoints

- **`/health`**: Comprehensive health check including database and cache
- **`/ready`**: Kubernetes readiness probe endpoint

### Structured Logging

All operations are logged with structured data:

```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "level": "info",
  "message": "HTTP request completed",
  "correlation_id": "req-123",
  "method": "GET",
  "path": "/api/v1/controls",
  "status": 200,
  "duration": "45ms"
}
```

### Audit Trail

Comprehensive audit logging for compliance:

```json
{
  "event_type": "audit",
  "action": "control_update",
  "user_id": "user-123",
  "resource_id": "control-456",
  "old_values": {"status": "draft"},
  "new_values": {"status": "active"},
  "correlation_id": "req-789"
}
```

## 🔧 Development

### Project Structure

```
internal/models/        # Domain models and entities
internal/services/      # Business logic interfaces
internal/repositories/  # Data access interfaces
pkg/                   # Reusable packages
```

### Adding New Features

1. **Define Models**: Add domain entities in `internal/models/`
2. **Create Interfaces**: Define service contracts in `internal/services/`
3. **Implement Services**: Add business logic implementations
4. **Add Repositories**: Implement data access layer
5. **Create Handlers**: Add HTTP/GraphQL endpoints
6. **Update Configuration**: Add any new config options

### Code Standards

- **Clean Architecture**: Maintain clear dependency boundaries
- **Interface-First**: Define interfaces before implementations
- **Comprehensive Documentation**: Document all public APIs
- **Error Handling**: Proper error propagation and logging
- **Testing**: Unit tests for business logic
- **Security**: Input validation and authorization

## 🚢 Deployment

### Environment Configuration

#### Development
```bash
GOEDU_APP_ENVIRONMENT=development
GOEDU_LOGGER_LEVEL=debug
```

#### Production
```bash
GOEDU_APP_ENVIRONMENT=production
GOEDU_LOGGER_LEVEL=info
GOEDU_AUTH_JWT_SECRET=secure-random-key
GOEDU_DATABASE_URI=mongodb://prod-cluster:27017
```

### Docker Deployment

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bin/server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/server .
EXPOSE 8080
CMD ["./server"]
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: goedu-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: goedu-backend
  template:
    metadata:
      labels:
        app: goedu-backend
    spec:
      containers:
      - name: backend
        image: goedu/backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: GOEDU_APP_ENVIRONMENT
          value: "production"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
```

## 📚 API Documentation

### Health Check
```http
GET /health
```

Response:
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "1.0.0",
  "checks": {
    "database": {"status": "healthy", "latency_ms": 12},
    "cache": {"status": "healthy", "latency_ms": 5}
  }
}
```

### Future API Endpoints

The application provides placeholders for:
- **GraphQL API**: `/api/v1/graphql`
- **REST Endpoints**: `/api/v1/*`
- **WebSocket**: For real-time notifications

## 🔒 Security

### Authentication & Authorization
- JWT-based authentication
- Role-based access control
- Multi-factor authentication support
- Session management

### Data Security
- Input validation and sanitization
- SQL injection protection
- XSS prevention
- CSRF protection

### Audit & Compliance
- Comprehensive audit logging
- Data retention policies
- Regulatory compliance support
- Tamper-evident logs

## 🤝 Contributing

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Make changes** following code standards
4. **Add tests** for new functionality
5. **Commit changes**: `git commit -m 'Add amazing feature'`
6. **Push to branch**: `git push origin feature/amazing-feature`
7. **Create Pull Request**

## 📄 License

This project is proprietary software. All rights reserved.

## 📞 Support

For support and questions:
- **Development Team**: [team@goedu.com]
- **Documentation**: See `/docs` folder
- **Issues**: Use GitHub Issues for bug reports

---

**Built with ❤️ for financial compliance professionals**