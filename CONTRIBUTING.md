# Contributing to AionMCP

Thank you for your interest in contributing to AionMCP! We welcome contributions from the community and want to make contributing as easy and transparent as possible.

---

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Setup](#development-setup)
- [Coding Guidelines](#coding-guidelines)
- [Commit Message Guidelines](#commit-message-guidelines)
- [Pull Request Process](#pull-request-process)
- [Reporting Bugs](#reporting-bugs)
- [Suggesting Enhancements](#suggesting-enhancements)
- [Community](#community)

---

## Code of Conduct

This project and everyone participating in it is governed by the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to [project-email@example.com].

---

## How Can I Contribute?

### 🐛 Reporting Bugs
- **Found a bug?** Check existing issues first
- Create a detailed bug report with:
  - Clear description of the bug
  - Steps to reproduce
  - Expected vs actual behavior
  - Environment details (OS, Go version, etc.)
  - Screenshots/logs if applicable

### 💡 Suggesting Enhancements
- Check existing issues/discussions
- Describe the feature clearly
- Explain why it would be useful
- Provide examples or mockups if applicable
- Consider the [Roadmap](./ROADMAP.md) to align with vision

### 📝 Documentation
- Improve existing docs
- Add examples
- Fix typos/unclear sections
- Add diagrams or visuals
- Create tutorials

### 🔧 Code Contributions
- Pick an issue labeled `good-first-issue`
- Work on features from [Roadmap](./ROADMAP.md)
- Add tests for new functionality
- Fix bugs with test coverage

### 🗣️ Community
- Answer questions in Discussions
- Share your use cases
- Provide feedback on features
- Help other contributors

---

## Development Setup

### Prerequisites
- **Go**: 1.21+
- **Git**: Latest version
- **Make**: For build commands (optional but recommended)
- **Docker**: For testing containerized deployments (optional)

### Local Development

1. **Fork & Clone**
   ```bash
   git clone https://github.com/YOUR_USERNAME/aionmcp.git
   cd aionmcp
   ```

2. **Create Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Install Dependencies**
   ```bash
   go mod download
   go mod verify
   ```

4. **Build**
   ```bash
   go build -o bin/aionmcp-server ./cmd/server
   ```

5. **Run Tests**
   ```bash
   go test ./...
   ```

6. **Run Server**
   ```bash
   ./bin/aionmcp-server
   ```

### Project Structure
```
aionmcp/
├── cmd/               # CLI entry points
│   └── server/        # Main server
├── internal/
│   ├── core/         # Core business logic
│   ├── autodocs/     # Auto-documentation engine
│   ├── selflearn/    # Learning system
│   └── adapters/     # API format adapters
├── pkg/
│   ├── importer/     # Spec import logic
│   ├── feedback/     # Feedback models
│   └── types/        # Shared types
├── examples/         # Example specs
├── docs/            # Documentation
└── tests/           # Integration tests
```

---

## Coding Guidelines

### Go Code Style
- **Format**: Use `gofmt` (enforced by CI)
  ```bash
  gofmt -w ./cmd ./internal ./pkg
  ```

- **Lint**: Use `golangci-lint`
  ```bash
  golangci-lint run ./...
  ```

- **Comments**: 
  - All exported functions must have doc comments
  - Use clear, complete sentences
  - Example:
    ```go
    // ImportOpenAPI loads an OpenAPI specification from the given path
    // and returns a list of Tool objects representing the API operations.
    func ImportOpenAPI(path string) ([]Tool, error) {
        // implementation
    }
    ```

### Error Handling
- Always return structured errors with context
- Use `fmt.Errorf` with `%w` for wrapping
- Log errors with sufficient context for learning system

Example:
```go
if err != nil {
    return fmt.Errorf("failed to parse spec %s: %w", specPath, err)
}
```

### Testing
- Minimum 70% code coverage
- Write unit tests for new functionality
- Include integration tests where applicable
- Use table-driven tests for multiple cases

Example:
```go
func TestImportOpenAPI(t *testing.T) {
    tests := []struct {
        name    string
        path    string
        wantErr bool
    }{
        {"valid spec", "specs/petstore.yaml", false},
        {"missing spec", "missing.yaml", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := ImportOpenAPI(tt.path)
            if (err != nil) != tt.wantErr {
                t.Errorf("got error %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Architecture Pattern
- Follow **Hexagonal/Clean Architecture**
- Separate concerns:
  - Core logic (business rules)
  - Adapters (external integrations)
  - Handlers (HTTP/API endpoints)
- Use interfaces for dependencies
- Make it testable and extensible

---

## Commit Message Guidelines

Use clear, descriptive commit messages following this format:

```
<type>: <subject>

<body>

<footer>
```

### Types
- **feat**: New feature
- **fix**: Bug fix
- **docs**: Documentation
- **test**: Tests
- **perf**: Performance improvement
- **refactor**: Code refactoring
- **chore**: Build, dependencies, config

### Examples
```
feat: add GraphQL schema validation

Implement JSON schema-based validation for GraphQL schemas
before importing. Validates type definitions and query structure.

Fixes #123
```

```
fix: handle timeout in API spec imports

Add configurable timeout for remote spec loading.
Defaults to 30s, prevents hanging on slow servers.

Closes #456
```

### Guidelines
- Use imperative mood: "add feature" not "added feature"
- Keep subject line < 50 characters
- Explain **what** and **why**, not **how**
- Reference related issues (#123)
- One logical change per commit

---

## Pull Request Process

### Before Submitting

1. **Update branch**
   ```bash
   git fetch origin
   git rebase origin/main
   ```

2. **Test locally**
   ```bash
   go test -v ./...
   go fmt ./...
   golangci-lint run ./...
   ```

3. **Update docs** if needed
   ```bash
   # Update CHANGELOG.md
   # Update relevant docs/
   ```

### Submitting PR

1. **Create descriptive PR title**
   - Follow commit message conventions
   - Example: "feat: add AsyncAPI schema validation"

2. **Write PR description**
   ```markdown
   ## Description
   Brief description of changes
   
   ## Type of Change
   - [ ] Bug fix
   - [ ] New feature
   - [ ] Documentation update
   
   ## Related Issues
   Fixes #123
   
   ## Testing
   - [ ] Unit tests added
   - [ ] Integration tests added
   - [ ] Manual testing done
   
   ## Checklist
   - [ ] Code follows style guidelines
   - [ ] Documentation updated
   - [ ] Tests pass locally
   - [ ] No breaking changes
   ```

3. **Request review** from maintainers

### After Submission
- Address review comments promptly
- Push updates to the same branch
- Don't force-push unless requested
- Keep PR focused and reasonably sized

---

## Reporting Bugs

### Good Bug Report Includes

1. **Summary**: One-line description
2. **Description**: Detailed explanation
3. **Steps to Reproduce**: Clear reproduction steps
4. **Expected Behavior**: What should happen
5. **Actual Behavior**: What actually happens
6. **Environment**:
   - OS and version
   - Go version (`go version`)
   - AionMCP version
   - Any relevant configurations
7. **Logs/Screenshots**: Relevant error messages or screenshots
8. **Additional Context**: Any other relevant info

### Example Bug Report
```markdown
## Summary
Server crashes when importing OpenAPI spec with circular references

## Description
When importing an OpenAPI specification containing circular schema 
references, the server crashes with a nil pointer error.

## Steps to Reproduce
1. Start server with `./aionmcp-server`
2. Import spec: `curl -X POST http://localhost:8080/api/v1/import-spec ...`
3. Use circular-ref-spec.yaml from examples/

## Expected Behavior
Server should handle circular references gracefully and skip or 
mark them appropriately.

## Actual Behavior
Server crashes with panic:
```
panic: runtime error: invalid memory address or nil pointer dereference
```

## Environment
- OS: macOS 14.1
- Go: 1.21
- AionMCP: v0.1.0
```

---

## Suggesting Enhancements

### Enhancement Proposal Includes

1. **Problem Statement**: What problem does it solve?
2. **Proposed Solution**: How should it work?
3. **Alternative Solutions**: Other approaches considered
4. **Use Cases**: Real-world scenarios
5. **Acceptance Criteria**: How to verify it works
6. **Additional Context**: Screenshots, mockups, examples

### Example Enhancement Proposal
```markdown
## Problem
When APIs return large datasets, tool execution takes too long 
and times out. Need way to handle pagination.

## Proposed Solution
Add automatic pagination support for JSON responses:
- Detect array responses > 1MB
- Split into pages based on configurable size
- Provide `next_page` token in response metadata
- Auto-fetch next page if requested

## Use Cases
- Large API result sets (1000+ items)
- Memory-constrained environments
- Streaming scenarios

## Acceptance Criteria
- [ ] Pagination support for OpenAPI operations
- [ ] Configurable page size
- [ ] Tests with 10k+ item responses
- [ ] Documentation and examples
```

---

## Development Workflow Example

1. **Check Issues**
   ```bash
   # Find something to work on
   # Look for: good-first-issue, help-wanted
   ```

2. **Create Branch**
   ```bash
   git checkout -b fix/import-timeout-issue
   ```

3. **Make Changes**
   ```bash
   # Edit files
   vim internal/importer/openapi.go
   
   # Test locally
   go test ./internal/importer -v
   ```

4. **Commit**
   ```bash
   git add .
   git commit -m "fix: add timeout for remote spec imports"
   ```

5. **Push & PR**
   ```bash
   git push origin fix/import-timeout-issue
   # Create PR on GitHub
   ```

6. **Address Feedback**
   ```bash
   # Make requested changes
   git add .
   git commit -m "refactor: improve error message clarity"
   git push origin fix/import-timeout-issue
   ```

7. **Merge**
   - Maintainer merges PR when approved
   - Branch is deleted

---

## Community

### Getting Help
- **Discussions**: https://github.com/kiransth77/aionmcp/discussions
- **Issues**: https://github.com/kiransth77/aionmcp/issues
- **Email**: [project-email@example.com]

### Resources
- **Documentation**: https://github.com/kiransth77/aionmcp/docs
- **Roadmap**: https://github.com/kiransth77/aionmcp/ROADMAP.md
- **Architecture**: https://github.com/kiransth77/aionmcp/docs/architecture.md

### Community Guidelines
- Be respectful and professional
- Help others when possible
- Share knowledge and ideas
- Celebrate contributions
- Report issues constructively

---

## Recognition

We recognize and thank all contributors! 

- Contributors are listed in [CONTRIBUTORS.md](CONTRIBUTORS.md)
- Major contributors are featured in releases
- Community highlights in monthly updates

---

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

## Questions?

Don't hesitate to ask! Open a discussion or check our [FAQ](./docs/FAQ.md).

Thank you for contributing to AionMCP! 🚀
