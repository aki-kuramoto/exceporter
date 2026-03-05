# Contributing to exceporter

Thank you for your interest in contributing!

## Reporting Issues

Please open a [GitHub Issue](https://github.com/aki-kuramoto/exceporter/issues) and include:

- Your OS and architecture
- Go version (`go version`)
- Steps to reproduce the problem
- Expected vs. actual behavior

## Submitting Changes

1. Fork the repository and create a branch from `main`.
2. Make your changes.
3. Run tests to make sure nothing is broken:
   ```sh
   go test ./...
   ```
4. Open a Pull Request with a clear description of what you changed and why.

## Development Setup

```sh
git clone https://github.com/aki-kuramoto/exceporter.git
cd exceporter
go mod download
go build ./...
```

See [docs/setup-gcloud-auth.md](docs/setup-gcloud-auth.md) for instructions on
setting up Google Drive authentication for local testing.

## Code Style

- Follow standard Go conventions (`go fmt`, `go vet`).
- Keep comments and error messages in English.
- Add or update tests for any logic changes in `internal/`.

## License

By contributing, you agree that your contributions will be licensed under the
[MIT License](LICENSE).
