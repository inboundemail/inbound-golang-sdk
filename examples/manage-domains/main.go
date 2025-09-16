package main

import (
	"context"
	"fmt"
	"log"
	"os"

	inbound "github.com/inboundemail/inbound-golang-sdk"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("INBOUND_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set the INBOUND_API_KEY environment variable")
	}

	// Create client
	client, err := inbound.NewClient(apiKey)
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}

	ctx := context.Background()

	// List existing domains
	fmt.Println("📋 Listing domains...")
	domainsResp, err := client.Domain().List(ctx, &inbound.GetDomainsRequest{
		Limit: inbound.Int(10),
	})

	if err != nil {
		log.Fatal("Failed to list domains:", err)
	}

	if domainsResp.Error != "" {
		log.Fatal("API Error:", domainsResp.Error)
	}

	fmt.Printf("Found %d domains:\n", len(domainsResp.Data.Data))
	for _, domain := range domainsResp.Data.Data {
		fmt.Printf("  • %s (Status: %s, Can Receive: %v)\n", 
			domain.Domain, domain.Status, domain.CanReceiveEmails)
	}

	// Example: Add a new domain (uncomment to use)
	/*
	fmt.Println("\n➕ Adding a new domain...")
	newDomainResp, err := client.Domain().Create(ctx, &inbound.PostDomainsRequest{
		Domain: "example.com",
	})

	if err != nil {
		log.Fatal("Failed to create domain:", err)
	}

	if newDomainResp.Error != "" {
		log.Fatal("API Error:", newDomainResp.Error)
	}

	fmt.Printf("✅ Domain created: %s\n", newDomainResp.Data.Domain)
	fmt.Printf("   Domain ID: %s\n", newDomainResp.Data.ID)
	fmt.Printf("   Status: %s\n", newDomainResp.Data.Status)
	
	// Get DNS records for the new domain
	fmt.Println("\n📋 DNS records required for verification:")
	recordsResp, err := client.Domain().GetDNSRecords(ctx, newDomainResp.Data.ID)
	if err != nil {
		log.Printf("Failed to get DNS records: %v", err)
	} else if recordsResp.Error != "" {
		log.Printf("API Error getting DNS records: %s", recordsResp.Error)
	} else {
		fmt.Printf("DNS records response: %+v\n", recordsResp.Data)
	}
	*/

	fmt.Println("\n✅ Domain management example completed!")
}
