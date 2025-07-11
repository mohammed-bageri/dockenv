# Contributing to dockenv

Thank you for your interest in contributing to dockenv! We welcome contributions from the community.

## Getting Started

### Prerequisites

- Go 1.24 or later
- Docker and Docker Compose (for testing)
- Git

### Development Setup

1. **Fork and clone the repository:**

   ```bash
   git clone https://github.com/mohammed-bageri/dockenv.git
   cd dockenv
   ```

2. **Install dependencies:**

   ```bash
   go mod download
   ```

3. **Build the project:**

   ```bash
   make build
   # or
   go build -o dockenv .
   ```

4. **Run tests:**
   ```bash
   make test
   # or
   go test ./...
   ```

## Development Workflow

### Making Changes

1. **Create a feature branch:**

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes:**

   - Follow Go best practices
   - Add tests for new functionality
   - Update documentation as needed

3. **Test your changes:**

   ```bash
   make check  # Runs fmt, vet, lint, and test
   ```

4. **Commit your changes:**

   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

   Follow [Conventional Commits](https://www.conventionalcommits.org/) format:

   - `feat:` for new features
   - `fix:` for bug fixes
   - `docs:` for documentation changes
   - `test:` for test additions/changes
   - `refactor:` for code refactoring

5. **Push and create a Pull Request:**
   ```bash
   git push origin feature/your-feature-name
   ```

### Code Style

- Use `gofmt` for formatting
- Follow Go conventions
- Write meaningful commit messages
- Add comments for complex logic
- Keep functions small and focused

### Testing

- Write unit tests for new functions
- Test CLI commands thoroughly
- Include integration tests when appropriate
- Ensure all tests pass before submitting PR

## Project Structure

```
dockenv/
├── cmd/                    # CLI commands
├── internal/
│   ├── config/            # Configuration management
│   ├── docker/            # Docker integration
│   ├── services/          # Service definitions
│   ├── systemd/           # Auto-start functionality
│   ├── templates/         # Docker Compose generation
│   └── utils/             # Utility functions
├── templates/             # Service templates
├── examples/              # Usage examples
└── .github/               # GitHub workflows
```

## Adding New Services

To add a new service (e.g., PostgreSQL):

1. **Add service definition in `internal/services/services.go`:**

   ```go
   "postgres": {
       Name:        "postgres",
       DisplayName: "PostgreSQL",
       Description: "PostgreSQL Database",
       DefaultPort: 5432,
       Template:    "postgres.yaml",
       Volumes:     []string{"postgres_data"},
       EnvVars: map[string]string{
           "DB_HOST": "127.0.0.1",
           "DB_PORT": "5432",
           // ...
       },
   },
   ```

2. **Create template in `templates/postgres.yaml`:**

   ```yaml
   postgres:
     image: postgres:15
     container_name: dockenv-postgres
     # ... rest of configuration
   ```

3. **Add embedded template in `internal/templates/templates.go`**

4. **Add tests and documentation**

## Adding New Commands

1. **Create command file in `cmd/`:**

   ```go
   package cmd

   import "github.com/spf13/cobra"

   var yourCmd = &cobra.Command{
       Use:   "your-command",
       Short: "Description",
       RunE:  runYourCommand,
   }

   func init() {
       rootCmd.AddCommand(yourCmd)
   }
   ```

2. **Add to root command and test thoroughly**

## Submitting Issues

### Bug Reports

Include:

- dockenv version (`dockenv --version`)
- Operating system and version
- Docker version (`docker --version`)
- Steps to reproduce
- Expected vs actual behavior
- Error messages/logs

### Feature Requests

Include:

- Clear description of the feature
- Use case and motivation
- Examples of how it would work
- Any alternatives considered

## Code Review Process

1. All submissions require review
2. CI tests must pass
3. Code coverage should not decrease significantly
4. Documentation must be updated for new features
5. Breaking changes require discussion

## Release Process

1. Update version in relevant files
2. Update CHANGELOG.md
3. Create release tag: `git tag v1.0.0`
4. GitHub Actions will build and release automatically

## Getting Help

- Check existing [issues](https://github.com/mohammed-bageri/dockenv/issues)
- Start a [discussion](https://github.com/mohammed-bageri/dockenv/discussions)
- Read the [documentation](README.md)

## Code of Conduct

Be respectful, inclusive, and professional in all interactions.

Thank you for contributing! 🚀
