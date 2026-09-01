# Changelog

## [1.0.0] - 2026-07-08

### Added
- Production-ready structured logging with log levels (DEBUG, INFO, WARN, ERROR)
- Comprehensive configuration system with .env and config.yaml support
- Centralized common response helpers for consistent API responses
- Complete test suite with mock repositories
- CI/CD pipeline with GitHub Actions
- Docker Compose and Kubernetes deployment support
- Nginx reverse proxy configuration
- PostgreSQL and Redis integration
- JWT authentication with refresh tokens
- WebSocket real-time updates
- Event-driven architecture with event bus
- Worker pool for async metric processing
- AI-powered insights and analysis
- Enterprise reporting (PDF, CSV, Excel)
- Remote file management
- Terminal command execution
- Docker and Kubernetes monitoring
- Security audit logging
- RBAC (Role-Based Access Control)

### Changed
- Replaced all `fmt.Println`/`log.Println` with structured logger
- Refactored config package to support all environment variables with defaults
- Updated database connection to use centralized config
- Updated Redis connection to use centralized config
- Improved error handling with consistent response format
- Enhanced graceful shutdown with proper resource cleanup

### Fixed
- Removed hardcoded values throughout the codebase
- Fixed duplicate code patterns in handlers
- Fixed inconsistent API response formats
- Fixed missing error handling in database operations

### Security
- JWT token expiration enforcement
- Password hashing with bcrypt
- API key authentication for agents
- Input validation on all endpoints
- Command allowlist support
- Audit logging for all sensitive operations
- CORS middleware configuration
- SQL injection protection via GORM

## [0.9.0] - 2026-06-15

### Added
- Initial enterprise release
- Machine enrollment and management
- Real-time metrics collection
- Docker container monitoring
- Kubernetes pod monitoring
- WebSocket-based live updates
- Basic authentication system
- Alert engine with configurable rules
- Dashboard with system overview
- API documentation portal