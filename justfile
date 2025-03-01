# GoPGO Justfile
# https://just.systems/

# List available commands
default:
    @just --list

# Build the application
build:
    go build -o gopgo

# Build with PGO profile
pgo-build PROFILE:
    go build -pgo={{PROFILE}} -o gopgo

# Run all benchmarks (both function-level and live server if available)
bench:
    go test -bench=. -benchmem

# Run specific benchmark test
bench-one BENCH:
    go test -bench={{BENCH}} -benchmem

# Analyze profile with pprof tool (requires graphviz for visualization)
# Example: just analyze-profile cpu.pprof
analyze-profile PROFILE:
    go tool pprof -http=:8081 {{PROFILE}}

# Full PGO workflow (instrument, generate load, rebuild)
pgo-workflow: build
    #!/bin/bash
    echo "🚀 Starting server in background..."
    ./gopgo &
    SERVER_PID=$!
    sleep 2

    echo "🧠 Starting CPU profile collection in background..."
    curl -o cpu.pprof "http://localhost:8080/debug/pprof/profile?seconds=35" > /dev/null 2>&1 &
    PROFILE_PID=$!

    echo "🔄 Generating load while profiling..."
    go test -bench=. -benchtime=4s

    echo "⏳ Waiting for profiling to complete..."
    wait $PROFILE_PID

    echo "🛑 Stopping server..."
    kill $SERVER_PID || true

    echo "✅ PGO workflow complete. Run benchmarks to compare performance."



# Compare benchmark performance before and after PGO (using live server)
benchmark-compare:
    #!/bin/bash
    echo "🔨 Building without PGO..."
    just build

    echo "🚀 Starting server in background for before benchmarks..."
    ./gopgo &
    SERVER_PID=$!
    sleep 2

    echo "🔍 Running live benchmarks before PGO..."
    go test -bench=. -benchmem -count=6 > benchmark-before.txt

    echo "🛑 Stopping server..."
    kill $SERVER_PID || true
    sleep 1

    echo "🔨 Building with PGO..."
    just pgo-build cpu.pprof

    echo "🚀 Starting server with PGO in background for after benchmarks..."
    ./gopgo &
    SERVER_PID=$!
    sleep 2

    echo "🔍 Running live benchmarks after PGO..."
    go test -bench=. -benchmem -count=6 > benchmark-after.txt

    echo "🛑 Stopping server..."
    kill $SERVER_PID || true

    echo "📊 Benchmark comparison:"
    benchstat benchmark-before.txt benchmark-after.txt
