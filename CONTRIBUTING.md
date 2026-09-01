# Contributing to InfraPilot Enterprise

Thank you for your interest in contributing to InfraPilot Enterprise! This guide will help you get started.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [How to Contribute](#how-to-contribute)
- [Development Setup](#development-setup)
- [Contribution Workflow](#contribution-workflow)
- [Issue Labels](#issue-labels)
- [Code Standards](#code-standards)
- [Testing](#testing)
- [Documentation](#documentation)
- [Community](#community)

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to:

- Be respectful and inclusive
- Welcome newcomers and help them get started
- Focus on constructive feedback
- Accept responsibility and apologize for mistakes

## Getting Started

### Prerequisites

- Go 1.24+
- Node.js 22+
- Docker & Docker Compose
- PostgreSQL 15+
- Redis 7+
- Git

### Quick Links

- 📖 [Documentation](https://github.com/venkaiswami/infrapilot-enterprise/tree/main/docs)
- 🐛 [Issue Tracker](https://github.com/venkaiswami/infrapilot-enterprise/issues)
- 💬 [Discussions](https://github.com/venkaiswami/infrapilot-enterprise/discussions)
- 🚀 [Plugin SDK](https://github.com/venkaiswami/infrapilot-enterprise/blob/main/docs/plugin-sdk.md)
- 🔌 [API Examples](https://github.com/venkaiswami/infrapilot-enterprise/blob/main/docs/api-sdk-examples.md)

## How to Contribute

There are many ways to contribute:

### 🐛 Report Bugs

- Check if the bug has already been reported in [Issues](https://github.com/venkaiswami/infrapilot-enterprise/issues)
- If not, create a new issue using the [Bug Report template](.github/ISSUE_TEMPLATE/bug_report.md)
- Include: steps to reproduce, expected behavior, actual behavior, environment details

### 💡 Suggest Features

- Check if the feature has been requested in [Issues](https://github.com/venkaiswami/infrapilot-enterprise/issues)
- If not, create a new issue using the [Feature Request template](.github/ISSUE_TEMPLATE/feature_request.md)
- Describe the problem you're solving, not just the solution

### 📝 Improve Documentation

- Fix typos, clarify confusing sections, add examples
- Documentation is as important as code!
- See [Documentation](#documentation) section below

### 🔌 Build Plugins

- Create plugins for cloud providers (AWS, Azure, GCP)
- Build notification integrations (Slack, Teams, Discord)
- See [Plugin SDK](https://github.com/venkaiswami/infrapilot-enterprise/blob/main/docs/plugin-sdk.md)

### 💻 Submit Code

- Fix bugs or implement features
- See [Contribution Workflow](#contribution-workflow) below

### 🎨 Design & UX

- Improve the frontend UI/UX
- Create themes or improve accessibility
- Design new dashboard widgets

## Development Setup

### 1. Fork and Clone

```bash
# Fork the repository on GitHub, then clone your fork
git clone https://github.com/YOUR_USERNAME/infrapilot-enterprise.git
cd infrapilot-enterprise

# Add upstream remote
git remote add upstream https://github.com/venkaiswami/infrapilot-enterprise.git
```

### 2. One-Click Installation

```bash
# Make the installer executable
chmod +x scripts/install.sh

# Run installation (installs everything automatically)
./scripts/install.sh --skip-deps

# Or use demo mode for testing
./scripts/install.sh --demo
```

### 3. Manual Setup (Alternative)

```bash
# Environment configuration
cp backend/.env.example backend/.env
# Edit backend/.env with your settings

# Start infrastructure
docker-compose up -d

# Backend setup
cd backend
go mod tidy
go run cmd/api/main.go

# Frontend setup (new terminal)
cd frontend
npm install
npm run dev

# Agent setup (optional, new terminal)
cd agent
go mod tidy
go run cmd/agent/main.go --server http://localhost:8080 --token YOUR_TOKEN
```

### 4. Verify Setup

```bash
# Backend health check
curl http://localhost:8080/healthz

# Frontend
open http://localhost:5173

# API test
curl http://localhost:8080/api/v1/health
```

## Contribution Workflow

### Step 1: Choose an Issue

- Look for issues labeled [`good first issue`](https://github.com/venkaiswami/infrapilot-enterprise/labels/good%20first%20issue) for beginners
- Look for [`help wanted`](https://github.com/venkaiswami/infrapilot-enterprise/labels/help%20wanted) for more complex tasks
- Comment on the issue to claim it
- A maintainer will assign it to you

### Step 2: Create a Branch

```bash
# Create a new branch from main
git checkout main
git pull upstream main
git checkout -b feature/your-feature-name

# Branch naming conventions:
# feature/add-cloudwatch-plugin
# fix/memory-leak-in-agent
# docs/improve-installation-guide
# refactor/simplify-event-bus
```

### Step 3: Make Changes

Follow our [Code Standards](#code-standards):

```bash
# Format Go code
go fmt ./...

# Run linting (if available)
golangci-lint run

# Run tests
go test ./...
```

### Step 4: Commit Changes

Use [Conventional Commits](https://www.conventionalcommits.org/):

```
# Format: <type>(<scope>): <subject>
# Types: feat, fix, docs, style, refactor, test, chore

feat(agent): add AWS CloudWatch plugin support
fix(backend): resolve memory leak in worker pool
docs(install): clarify Docker Compose setup steps
```

```bash
git add .
git commit -m "feat(agent): add AWS CloudWatch plugin"
```

### Step 5: Push and Create PR

```bash
# Push to your fork
git push origin feature/your-feature-name

# Create a Pull Request on GitHub
# Link the issue: "Closes #123"
# Fill out the PR template completely
```

### Step 6: Code Review

- A maintainer will review your PR
- Address any requested changes
- Once approved, your PR will be merged

## Issue Labels

We use labels to organize issues:

| Label | Description | For Beginners |
|-------|-------------|---------------|
| `good first issue` | Beginner-friendly tasks | ✅ Yes |
| `help wanted` | Extra attention needed | ✅ Yes |
| `documentation` | Documentation improvements | ✅ Yes |
| `bug` | Something isn't working | ⚠️ Intermediate |
| `enhancement` | New feature or request | ⚠️ Intermediate |
| `plugin` | Plugin development | ⚠️ Intermediate |
| `security` | Security-related issues | 🔴 Advanced |
| `infrastructure` | DevOps/infrastructure | 🔴 Advanced |
| `ci/cd` | CI/CD improvements | ⚠️ Intermediate |
| `testing` | Test coverage | ✅ Yes |

## Code Standards

### Go (Backend & Agent)

```go
// Use standard Go formatting
go fmt ./...

// Naming conventions:
// - Use camelCase for variables and functions
// - Use PascalCase for exported types and functions
// - Use descriptive names: userService, not us

// Error handling:
func doSomething() error {
    if err := riskyOperation(); err != nil {
        return fmt.Errorf("failed to do something: %w", err)
    }
    return nil
}

// Comments for exported functions:
// GetUser retrieves a user by ID.
// Returns ErrUserNotFound if the user doesn't exist.
func GetUser(id string) (*User, error) {
    // ...
}

// Use context for cancellation:
func (s *Service) FetchData(ctx context.Context, url string) ([]byte, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    // ...
}
```

### JavaScript/React (Frontend)

```javascript
// Use functional components with hooks
import React, { useState, useEffect } from 'react';

// Component naming: PascalCase
const MachineCard = ({ machine }) => {
  // Hooks at the top
  const [loading, setLoading] = useState(false);

  // Event handlers: handle<Event>
  const handleClick = () => {
    console.log('Clicked');
  };

  return (
    <div className="machine-card">
      <h3>{machine.name}</h3>
    </div>
  );
};

// Use descriptive variable names
const fetchMachines = async () => {
  // Good
  const response = await api.machines.list();
  const machines = response.data;

  // Avoid
  const res = await api.machines.list();
  const m = res.data;
};
```

### General Guidelines

- **DRY** - Don't Repeat Yourself
- **KISS** - Keep It Simple, Stupid
- **YAGNI** - You Aren't Gonna Need It
- Write self-documenting code
- Add comments for complex logic, not obvious code
- Keep functions small and focused
- Use meaningful names

## Testing

### Backend Tests

```bash
cd backend

# Run all tests
go test -v -count=1 ./...

# Run specific test
go test -v -run TestCreateMachine

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Frontend Tests

```bash
cd frontend

# Run tests
npm run test

# Run with coverage
npm run test -- --coverage
```

### Write Tests

```go
// Backend: test file naming: foo_test.go
package handlers

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestCreateMachine(t *testing.T) {
    // Setup
    handler := NewMachineHandler()

    // Test
    result, err := handler.Create(machineData)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "web-server", result.Name)
}
```

```javascript
// Frontend: Component.test.jsx
import { render, screen } from '@testing-library/react';
import { MachineCard } from './MachineCard';

test('renders machine name', () => {
  const machine = { name: 'web-server', status: 'online' };
  render(<MachineCard machine={machine} />);
  
  expect(screen.getByText('web-server')).toBeInTheDocument();
});
```

## Documentation

### Updating Documentation

1. Edit `.md` files in `docs/`
2. Follow Markdown best practices
3. Keep it clear and concise
4. Add examples where helpful

### Documentation Structure

```
docs/
├── README.md              # Main documentation hub
├── architecture.md        # System architecture
├── api.md                 # API reference
├── plugin-sdk.md          # Plugin development guide
├── api-sdk-examples.md    # SDK examples
├── installation.md        # Installation guide
├── deployment.md          # Deployment guide
├── agent.md               # Agent documentation
├── backend.md             # Backend development
├── frontend.md            # Frontend development
└── troubleshooting.md     # Common issues
```

### Documentation Checklist

- [ ] Clear and concise language
- [ ] Code examples tested and working
- [ ] Links are valid
- [ ] Formatting is consistent
- [ ] Screenshots added (if UI changes)

## Community

### Getting Help

- 💬 [GitHub Discussions](https://github.com/venkaiswami/infrapilot-enterprise/discussions) - Ask questions
- 🐛 [GitHub Issues](https://github.com/venkaiswami/infrapilot-enterprise/issues) - Report bugs
- 📧 Email: venkaiswami@pm.me

### Stay Updated

- ⭐ Star the repository
- 👀 Watch for releases
- 📢 Follow on [LinkedIn](https://linkedin.com/in/venkaiswami)

### Recognition

Contributors are recognized in:
- Release notes
- README.md contributors section
- Project website

### Plugin Ecosystem

Share your plugins:
1. Create a repository with your plugin
2. Document it following our SDK guide
3. Add to the [Plugin SDK](docs/plugin-sdk.md) examples
4. Share in [Discussions](https://github.com/venkaiswami/infrapilot-enterprise/discussions)

### Mentorship

New to open source? We're here to help!

- Comment on a [`good first issue`](https://github.com/venkaiswami/infrapilot-enterprise/labels/good%20first%20issue)
- We'll pair you with a mentor
- Get guidance through your first contribution

## Frequently Asked Questions

**Q: Do I need to create a plugin to contribute?**
A: No! We welcome all contributions: bug fixes, documentation, tests, features.

**Q: Can I contribute plugins for cloud providers?**
A: Absolutely! See the [Plugin SDK](docs/plugin-sdk.md) guide.

**Q: How long does PR review take?**
A: Typically 2-3 business days. Complex PRs may take longer.

**Q: What if my PR is rejected?**
A: That's normal! We'll provide feedback and help you improve it.

**Q: Can I work on multiple issues?**
A: Yes, but please claim them first so others know you're working on them.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to InfraPilot Enterprise! 🚀

**Questions?** Reach out in [Discussions](https://github.com/venkaiswami/infrapilot-enterprise/discussions) or email venkaiswami@pm.me