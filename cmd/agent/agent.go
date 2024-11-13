package main

import (
	"flag"
	"log"

	"ddgo/agent"
)

func main() {
	// command-line flag for the server URL with a default value of "http://localhost:8080".
	serverURL := flag.String("server", "http://localhost:8080", "URL of the central metrics server")

	// parse the command-line flags.
	flag.Parse()

	// create a new agent instance with the specified server URL.
	a, err := agent.NewAgent(*serverURL)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	// start the agent and log a fatal error if it fails to start.
	log.Fatal(a.Start())
}
