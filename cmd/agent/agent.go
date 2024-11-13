package main

import (
	"flag"
	"log"

	"ddgo/agent"
)

func main() {
	serverURL := flag.String("server", "http://localhost:8080", "URL of the central metrics server")
	flag.Parse()

	a, err := agent.NewAgent(*serverURL) // start new server
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	log.Fatal(a.Start())
}
