# GoPGO

GoPGO is a demonstration project for Go's Profile-Guided Optimization (PGO) feature introduced in Go 1.20. It provides a simple web service with endpoints that showcase the effects of PGO on different code paths.

## Overview

The project implements a HTTP server with the following endpoints:

- `/status` - A simple status check endpoint
- `/compute/:complexity` - A CPU-intensive endpoint that performs SHA-256 hashing with adjustable complexity
- `/debug/pprof/*` - Standard Go profiling endpoints

## Features

- Demonstrates Go's PGO workflow with a real-world web service example
- Includes benchmarks to measure performance before and after applying PGO
- Provides automation scripts via Justfile to simplify the PGO workflow
- Shows the impact of PGO on both hot and cold execution paths

## Requirements

- Go 1.20 or later (for PGO support)
- [just](https://github.com/casey/just) command runner (for workflow automation)
- [benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat) (for benchmark comparison)

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/GoPGO.git
cd GoPGO

# Build the project
just build
# Or without just:
go build -o gopgo
```

## Usage

### Running the server

```bash
# Using just
just run

# Or directly
./gopgo
```

This starts the server on port 8080.

### Available endpoints

- `http://localhost:8080/status` - Returns server status
- `http://localhost:8080/compute/10` - Performs 10 iterations of hashing
- `http://localhost:8080/compute/100` - Performs 100 iterations of hashing
- `http://localhost:8080/compute/1000` - Performs 1000 iterations of hashing (complexity capped at 1000 for safety)
- Add `?cold=true` to any compute endpoint to also execute the cold path

### PGO workflow

The complete PGO workflow is automated using the Justfile:

```bash
# Run the complete PGO workflow: build, profile, rebuild with PGO
just pgo-workflow

# Compare performance before and after PGO
just benchmark-compare
```

### Additional commands

```bash
# Run all benchmarks
just bench

# Run a specific benchmark
just bench-one BenchmarkStatusEndpointLive

# Format code
just fmt

# Run linter
just lint

# Run tests
just test
```

## Project Structure

- `main.go` - Web server implementation with compute and status endpoints
- `main_test.go` - Benchmarks for measuring performance
- `justfile` - Automation scripts for the PGO workflow
- `CLAUDE.md` - Development guidelines and commands

## Benchmarking Results

The project includes benchmarks that measure the performance of various endpoints. The results show the impact of PGO on different types of workloads:

- Status endpoint (simple): ~4.7% improvement with PGO
- Compute endpoint (CPU-intensive): ~26-29% performance regression with PGO
- Cold path (rarely executed): No significant change with PGO

## Learn More

To learn more about Go's PGO implementation:

- [Go PGO Documentation](https://go.dev/doc/pgo)
- [Go 1.20 Release Notes](https://go.dev/doc/go1.20)

## License

[MIT License](LICENSE)