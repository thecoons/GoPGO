package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

// Define a flag to specify the server address
var serverAddr = flag.String("server", "http://localhost:8080", "Server address for benchmarks")

// BenchmarkStatusEndpointLive benchmarks the /status endpoint on a live server
func BenchmarkStatusEndpointLive(b *testing.B) {
	// Skip if running in short mode (-short flag)
	if testing.Short() {
		b.Skip("Skipping live server test in short mode")
	}

	// Check if server is running
	_, err := http.Get(*serverAddr + "/status")
	if err != nil {
		b.Fatalf("Server not available at %s: %v - Start server first using 'just serve-profile'", *serverAddr, err)
	}

	// Create a client with timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Reset the timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		resp, err := client.Get(*serverAddr + "/status")
		if err != nil {
			b.Fatalf("Request failed: %v", err)
		}
		_, err = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		if err != nil {
			b.Fatalf("Error reading response: %v", err)
		}
	}
}

// BenchmarkComputeEndpointLive benchmarks the /compute endpoint with different complexity levels
// against a live server
func BenchmarkComputeEndpointLive1000(b *testing.B) {
	benchmarkComputeEndpointLive(b, 1000)
}

func BenchmarkComputeEndpointLive10000(b *testing.B) {
	benchmarkComputeEndpointLive(b, 10000)
}

func BenchmarkComputeEndpointLive100000(b *testing.B) {
	benchmarkComputeEndpointLive(b, 100000)
}

// benchmarkComputeEndpointLive is a helper function to benchmark the compute endpoint with a given complexity
func benchmarkComputeEndpointLive(b *testing.B, complexity int) {
	// Skip if running in short mode (-short flag)
	if testing.Short() {
		b.Skip("Skipping live server test in short mode")
	}

	// Check if server is running
	_, err := http.Get(*serverAddr + "/status")
	if err != nil {
		b.Fatalf("Server not available at %s: %v - Start server first using 'just serve-profile'", *serverAddr, err)
	}

	// Create a client with timeout
	client := &http.Client{
		Timeout: 15 * time.Second, // Longer timeout for compute-intensive operations
	}

	url := fmt.Sprintf("%s/compute/%d?data=benchmark-data", *serverAddr, complexity)

	// Reset the timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		resp, err := client.Get(url)
		if err != nil {
			b.Fatalf("Request failed: %v", err)
		}
		_, err = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		if err != nil {
			b.Fatalf("Error reading response: %v", err)
		}
	}
}

// BenchmarkColdPathLive benchmarks the cold path against a live server
func BenchmarkColdPathLive(b *testing.B) {
	// Skip if running in short mode (-short flag)
	if testing.Short() {
		b.Skip("Skipping live server test in short mode")
	}

	// Check if server is running
	_, err := http.Get(*serverAddr + "/status")
	if err != nil {
		b.Fatalf("Server not available at %s: %v - Start server first using 'just serve-profile'", *serverAddr, err)
	}

	// Create a client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	url := fmt.Sprintf("%s/compute/10?data=benchmark-data-with-multiple-parts&cold=true", *serverAddr)

	// Reset the timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		resp, err := client.Get(url)
		if err != nil {
			b.Fatalf("Request failed: %v", err)
		}
		_, err = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		if err != nil {
			b.Fatalf("Error reading response: %v", err)
		}
	}
}
