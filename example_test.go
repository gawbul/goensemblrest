package goensemblrest_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	ensembl "github.com/gawbul/goensemblrest"
)

func ExampleNewClient() {
	// Create a new Ensembl REST client with default settings (15 req/s, 60s timeout).
	client, err := ensembl.NewClient()
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	fmt.Println("Base URL:", client.BaseURL())
	// Output:
	// Base URL: https://rest.ensembl.org
}

func ExampleClient_GetLookupByID() {
	client, err := ensembl.NewClient(
		ensembl.WithTimeout(30 * time.Second),
	)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var gene ensembl.LookupRecord
	err = client.GetLookupByID(ctx, "ENSG00000157764", &gene, ensembl.WithQuery("expand", "1"))
	if err != nil {
		// Handle potential errors
		if errors.Is(err, ensembl.ErrNotFound) {
			fmt.Println("Gene not found")
			return
		}
		log.Printf("request error: %v", err)
		return
	}

	fmt.Println("Lookup request prepared successfully")
}

func ExampleClient_Call() {
	client, err := ensembl.NewClient()
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	// Dynamic endpoint dispatch by name
	var result map[string]any
	err = client.Call(
		context.Background(),
		"getInfoPing",
		nil,
		nil,
		&result,
	)
	if err != nil {
		log.Printf("ping failed: %v", err)
		return
	}
}
