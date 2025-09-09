package main

import (
	"context"
	"fmt"
	"log"
	"os"

	inbound "github.com/R44VC0RP/inbound-golang-sdk"
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

	// Send a simple email
	fmt.Println("Sending email...")
	resp, err := client.Email().Send(ctx, &inbound.PostEmailsRequest{
		From:    "hello@yourdomain.com",
		To:      "recipient@example.com",
		Subject: "Hello from Inbound Go SDK!",
		Text:    inbound.String("This is a test email sent using the Inbound Go SDK."),
		HTML:    inbound.String("<h1>Hello!</h1><p>This is a test email sent using the <strong>Inbound Go SDK</strong>.</p>"),
		Tags: []inbound.EmailTag{
			{Name: "environment", Value: "example"},
			{Name: "sdk", Value: "go"},
		},
	}, nil)

	if err != nil {
		log.Fatal("Failed to send email:", err)
	}

	if resp.Error != "" {
		log.Fatal("API Error:", resp.Error)
	}

	fmt.Printf("✅ Email sent successfully!\n")
	fmt.Printf("   Email ID: %s\n", resp.Data.ID)
	if resp.Data.MessageID != nil {
		fmt.Printf("   Message ID: %s\n", *resp.Data.MessageID)
	}
	if resp.Data.Status != nil {
		fmt.Printf("   Status: %s\n", *resp.Data.Status)
	}
}
