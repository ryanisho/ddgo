package main

import (
	"flag"
	"log"

	"ddgo/agent"
)

// main is the entry point of the program.
func main() {
	// Define a command-line flag for the server URL with a default value of "http://localhost:8080".
	serverURL := flag.String("server", "http://localhost:8080", "URL of the central metrics server")

	// Parse the command-line flags.
	flag.Parse()

	// Create a new agent instance with the specified server URL.
	a, err := agent.NewAgent(*serverURL)
	if err != nil {
		// Log a fatal error and exit if the agent creation fails.
		log.Fatalf("Failed to create agent: %v", err)
	}

	// Start the agent and log a fatal error if it fails to start.
	log.Fatal(a.Start())
}
