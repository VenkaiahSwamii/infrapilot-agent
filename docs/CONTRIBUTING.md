# Contributing to InfraPilot Enterprise

We love your input! We want to make contributing to InfraPilot as easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Becoming a maintainer

## Development Process

We use GitHub to host code, to track issues and feature requests, as well as accept pull requests.

### 1. Fork & Clone

```bash
git clone https://github.com/your-username/infrapilot-enterprise.git
cd infrapilot-enterprise
git remote add upstream https://github.com/venkaiswami/infrapilot-enterprise.git
```

### 2. Create a Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

### 3. Set Up Development Environment

**Prerequisites:**
- Go 1.24+
- Node.js 22+
- Docker & Docker Compose
- PostgreSQL 15+ (for full integration testing)

**Start infrastructure:**
```bash
docker-compose up -d postgres redis
```

**Start backend (terminal 1):**
```bash
cd backend
cp .env.example .env
go mod tidy
go run cmd/api/main.go
```

**Start frontend (terminal 2):**
```bash
cd frontend
npm install
npm run dev
```

**Run agent (terminal 3):**
```bash
cd agent
go run cmd/agent/main.go --server http://localhost:8080 --token YOUR_TOKEN
```

### 4. Make Changes

- Follow the [code style](#code-style) guidelines
- Write or update tests as needed
- Update documentation for API changes
- Ensure all tests pass before committing

### 5. Test Your Changes

```bash
# Backend tests
cd backend
go test -v -count=1 ./...

# Agent tests
cd agent
go test -v -count=1 ./...

# Frontend build check
cd frontend
npm run build

# Full Docker build
docker-compose build
```

### 6. Commit & Push

```bash
git add .
git commit -m "feat: clear description of your change"
git push origin feature/your-feature-name
```

### 7. Open a Pull Request

1. Go to the original repository on GitHub
2. Click "New Pull Request"
3. Select your branch
4. Fill out the PR template with:
   - What changes you made
   - Why you made them
   - How to test them
   - Screenshots (for UI changes)

## Code Style

### Go (Backend & Agent)

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` or `go fmt` before committing
- Run `golangci-lint` for additional checks
- Use meaningful variable names; avoid single-letter names except in short scopes
- Handle all errors; never use `_` to discard errors without reason
- Write godoc comments for all exported types and functions

```go
// Bad
func GetMachine(mid string) (*Machine, error)

// Good
// GetMachineByID retrieves a machine by its UUID.
func GetMachineByID(ctx context.Context, machineID uuid.UUID) (*Machine, error)
```

### JavaScript / React (Frontend)

- Use functional components with hooks
- No class components unless absolutely necessary
- Use meaningful component and file names
- Keep components focused and under 300 lines where possible
- Extract reusable logic into custom hooks

```jsx
// Bad
function Page() { /* 500 lines */ }

// Good
function EnterpriseDashboard() {
  const machines = useMachines();
  const alerts = useAlerts();
  return (
    <div>
      <Header />
      <KpiCards machines={machines} alerts={alerts} />
      <MachinesTable machines={machines} />
    </div>
  );
}
```

### Commit Messages

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add alert filtering by severity
^    ^
|    └─> summary in present tense (lowercase, no period)
|
└─────> type: feat, fix, docs, style, refactor, test, chore, ci

fix: resolve WebSocket reconnection loop
docs: update API documentation
test: add rate limiter unit tests
refactor: extract metric validation to service layer
```

## Testing Guidelines

### Unit Tests
- Test one thing per test case
- Use descriptive test names: `TestFunctionName_Scenario`
- Use table-driven tests for multiple scenarios
- Mock external dependencies (database, Redis, network)

### Integration Tests
- Tag integration tests with `//go:build integration`
- Use test containers for PostgreSQL and Redis
- Test full request/response cycles

### Test Coverage Goals
| Layer | Target |
|-------|--------|
| Domain | 90%+ |
| Services | 80%+ |
| Handlers | 70%+ |
| Infrastructure | 60%+ |

## Pull Request Guidelines

- Keep PRs focused on a single concern
- Include tests for new functionality
- Update documentation for changed APIs
- Ensure CI pipeline passes all checks
- Request review from at least one maintainer

### PR Size Guidelines
- **Small PR (< 200 lines):** Quick review, ideal for bug fixes
- **Medium PR (200-500 lines):** Normal feature work
- **Large PR (> 500 lines):** Should be broken into smaller PRs if possible

## Review Process

1. **Author** opens PR with clear description
2. **CI** runs automated checks (lint, test, build, security)
3. **Reviewer(s)** provide feedback:
   - Code quality and style
   - Test coverage
   - Documentation
   - Security implications
4. **Author** addresses feedback with additional commits
5. **Approval** from at least one maintainer
6. **Merge** by a maintainer (squash or rebase)

## Project Structure

```
infrapilot-enterprise/
├── agent/           # Go agent - runs on monitored servers
│   ├── cmd/         # Agent entry point
│   └── internal/    # Agent internals (collectors, senders)
├── backend/         # Go API server
│   ├── cmd/         # Server, worker, scheduler entry points
│   └── internal/    # Backend internals
├── frontend/        # React dashboard
│   └── src/         # Frontend source
├── docs/            # Documentation
├── k8s/             # Kubernetes manifests
├── grafana/         # Grafana dashboards
├── scripts/         # Deployment scripts
├── installer/       # Windows/Linux installers
└── website/         # Marketing website
```

## Questions?

- Open a [GitHub Discussion](https://github.com/venkaiswami/infrapilot-enterprise/discussions)
- Check the [documentation](https://github.com/venkaiswami/infrapilot-enterprise/tree/main/docs)
- Email: venkaiswami@pm.me

## License

By contributing, you agree that your contributions will be licensed under the MIT License.