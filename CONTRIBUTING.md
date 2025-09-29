# Contributing to ICS Terraform Provider

Thank you for your interest in contributing to the ICS Terraform Provider! This guide will help you get started with development and testing.

## Development Setup

### Prerequisites

- **Go**: Version 1.21 or later
- **Terraform**: Version 1.0 or later
- **Make**: For using the provided Makefile
- **ICS API Token**: For testing (obtain from ICS dashboard)

### Building the Provider

To build the provider binary:

```bash
go build -o terraform-provider-ics
```

### Local Installation

For local development and testing, install the provider to your local Terraform plugins directory:

```bash
make install
```

This will:
1. Build the provider binary
2. Install it to `~/.terraform.d/plugins/local.dev/UK2Group/ics/0.1/darwin_arm64/`

### Testing

#### Unit Tests

Run the unit tests:

```bash
go test ./...
```

#### Acceptance Tests

Run acceptance tests (requires valid ICS API credentials):

```bash
TF_ACC=1 go test ./... -v -timeout 120m
```

**Note**: Acceptance tests will create real resources and may incur charges.

### Code Quality

#### Formatting

Format your code:

```bash
go fmt ./...
```

#### Linting

Run the linter:

```bash
golangci-lint run
```

### Documentation

#### Generate Documentation

Update the documentation using tfplugindocs:

```bash
go generate ./...
```

This will update the documentation in the `docs/` directory based on the provider schema.

## Development Workflow

### 1. Making Changes

1. **Fork the repository** and clone your fork
2. **Create a feature branch** from main
3. **Make your changes** with appropriate tests
4. **Run tests** to ensure everything works
5. **Update documentation** if needed

### 2. Testing Your Changes

1. **Build and install** the provider locally:
   ```bash
   make install
   ```

2. **Test with the examples**:
   ```bash
   cd example/
   terraform init
   terraform plan
   ```

3. **Run the test suite**:
   ```bash
   go test ./...
   ```

### 3. Submitting Changes

1. **Push your changes** to your fork
2. **Create a pull request** with:
   - Clear description of changes
   - Any breaking changes noted
   - Test results
   - Updated documentation

## Project Structure

```
├── docs/                    # Generated documentation
├── example/                 # Example configurations
├── internal/
│   └── provider/           # Provider implementation
├── .github/workflows/      # CI/CD workflows
├── .goreleaser.yml        # Release configuration
├── main.go                # Provider entry point
├── Makefile              # Build tasks
└── README.md             # User documentation
```

## Release Process

Releases are automated via GitHub Actions when tags are pushed:

1. **Create a tag**: `git tag v1.0.0`
2. **Push the tag**: `git push origin v1.0.0`
3. **GitHub Actions** will automatically:
   - Build binaries for all platforms
   - Create a GitHub release
   - Sign the release with GPG

## Code Style Guidelines

- Follow standard Go formatting (`go fmt`)
- Write comprehensive tests for new features
- Update documentation for any user-facing changes
- Use clear, descriptive commit messages
- Keep functions focused and well-documented

## Getting Help

- **Issues**: Report bugs or request features via GitHub Issues
- **Discussions**: Use GitHub Discussions for questions
- **ICS API**: Consult ICS documentation for API-related questions

## License

By contributing, you agree that your contributions will be licensed under the same license as the project (Mozilla Public License 2.0).