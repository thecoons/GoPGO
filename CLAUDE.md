# GoPGO Development Guidelines

## Build Commands
- Project uses Justfile for common operations: `just <command>`
- Build: `just build` or `go build -o gopgo`
- Run: `just run` or `go run main.go`
- Format: `just fmt`
- Lint: `just lint`
- Test: `just test` or `go test ./...`
- Test single function: `just test-one TestFunctionName`
- Benchmark: `just bench`
- PGO workflow: `just pgo-workflow` (instrument, load, rebuild)

## Code Style Guidelines
- Format code with: `gofmt -s -w .` or `go fmt ./...`
- Lint with: `golint ./...` or `staticcheck ./...`
- Imports: Group standard library, third-party, and internal imports
- Error handling: Always check and handle errors, prefer early returns
- Naming: Use camelCase for variables, PascalCase for exported symbols
- Functions: Keep functions short and focused on a single responsibility
- Comments: Use godoc style comments for all exported symbols
- Types: Use strong typing, avoid interface{} when possible
- Avoid globals and use dependency injection instead

For PGO workflow:
1. Build with instrumentation
2. Generate representative load
3. Rebuild with profile data
4. Benchmark to confirm improvements