package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	// Wait a moment for server to start
	time.Sleep(2 * time.Second)

	// Test health endpoint
	fmt.Println("Testing health endpoint...")
	resp, err := http.Get("http://localhost:8080/api/v1/health")
	if err != nil {
		fmt.Printf("Error calling health endpoint: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	var healthResp map[string]interface{}
	if err := json.Unmarshal(body, &healthResp); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	fmt.Printf("Health Response: %+v\n", healthResp)

	// Test tools endpoint
	fmt.Println("\nTesting tools endpoint...")
	resp, err = http.Get("http://localhost:8080/api/v1/mcp/tools")
	if err != nil {
		fmt.Printf("Error calling tools endpoint: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	var toolsResp map[string]interface{}
	if err := json.Unmarshal(body, &toolsResp); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	fmt.Printf("Tools Response: %+v\n", toolsResp)
}