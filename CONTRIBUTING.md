# Contributing

## Development

Requirements:
- Go 1.18+

Useful commands:

```bash
go test ./...
gofmt -w .
```

## Guidelines

- Keep public models explicit and documentation-driven.
- Do not replace typed request or response fields with raw JSON containers.
- Prefer small, reviewable pull requests.
- Add or update tests for transport behavior and response decoding when changing models.

## Pull Requests

Please include:
- what changed;
- why it changed;
- whether the change is based on a documentation update;
- tests added or updated.
