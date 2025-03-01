package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"strconv"
	"strings"
	"time"
)

// This example demonstrates a simple web service with two endpoints:
// 1. /compute/:complexity - CPU intensive operation with varying complexity
// 2. /status - Simple status check

// computeIntensiveHash performs a CPU-intensive operation
// The complexity parameter determines how many iterations of hashing are performed
func computeIntensiveHash(data string, complexity int) string {
	result := []byte(data)

	// This is our "hot path" that will benefit most from PGO
	for i := 0; i < complexity; i++ {
		h := sha256.New()
		h.Write(result)
		result = h.Sum(nil)
	}

	return hex.EncodeToString(result)
}

// coldPath is a function that is rarely called in our workload
// PGO should deprioritize optimizing this
func coldPath(input string) string {
	parts := strings.Split(input, "-")
	var result string

	// Some complex logic that's rarely executed
	for i := len(parts) - 1; i >= 0; i-- {
		result += strings.ToUpper(parts[i])
		if i > 0 {
			result += "_"
		}
	}

	return result
}

func main() {
	// Register pprof handlers
	// Register just the index handler and let it handle the other endpoints
	http.HandleFunc("/debug/pprof/", pprof.Index)

	// Status endpoint - simple and fast
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server is running. Time: %s", time.Now().Format(time.RFC3339))
	})

	// Compute endpoint - CPU intensive
	http.HandleFunc("/compute/", func(w http.ResponseWriter, r *http.Request) {
		// Extract complexity from URL path
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 3 {
			http.Error(w, "Invalid path. Use /compute/:complexity", http.StatusBadRequest)
			return
		}

		complexityStr := pathParts[2]
		complexity, err := strconv.Atoi(complexityStr)
		if err != nil {
			http.Error(w, "Complexity must be an integer", http.StatusBadRequest)
			return
		}

		// Cap complexity to prevent DOS
		if complexity > 1000 {
			complexity = 1000
		}

		// Get input data from query parameter or use default
		inputData := r.URL.Query().Get("data")
		if inputData == "" {
			inputData = "default-input-data"
		}

		// Perform the computation
		startTime := time.Now()
		result := computeIntensiveHash(inputData, complexity)
		duration := time.Since(startTime)

		// Only call the cold path occasionally
		var coldResult string
		if r.URL.Query().Get("cold") == "true" {
			coldResult = coldPath(inputData)
		}

		// Return the result and timing information
		response := fmt.Sprintf("Input: %s\nComplexity: %d\nResult: %s\nDuration: %s\n",
			inputData, complexity, result, duration)

		if coldResult != "" {
			response += fmt.Sprintf("Cold path result: %s\n", coldResult)
		}

		fmt.Fprint(w, response)
	})

	// Start the server
	port := 8080
	fmt.Printf("Starting server on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
