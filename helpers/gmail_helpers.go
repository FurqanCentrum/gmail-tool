package helpers

import (
	"context"
	"fmt"
	"gmail_tool_project/tools"
	"log"
	"strings"
)

// InitializeGmailTool initializes and returns an instance of GmailTool
func InitializeGmailTool(ctx context.Context) *tools.GmailTool {
	gmailTool, err := tools.InitializeGmailTool(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize Gmail tool: %v", err)
	}
	return gmailTool
}

// DisplayEmailIDs lists and displays email IDs based on a query
func DisplayEmailIDs(ctx context.Context, gmailTool *tools.GmailTool, query string) {
	emailIDs, err := gmailTool.ListEmailIDs(ctx, query)
	if err != nil {
		log.Fatalf("Failed to list email IDs: %v", err)
	}

	fmt.Println("Email IDs:")
	for _, id := range emailIDs {
		fmt.Println(id)
	}
}

// DisplayEmailDetails lists and displays email details based on a query
func DisplayEmailDetails(ctx context.Context, gmailTool *tools.GmailTool, query string) {
	emailIDs, err := gmailTool.ListEmailIDs(ctx, query)
	if err != nil {
		log.Fatalf("Failed to list email IDs: %v", err)
	}

	fmt.Println("Email Details:")
	for _, id := range emailIDs {
		msg, err := gmailTool.GetEmailDetails(ctx, id)
		if err != nil {
			log.Printf("Failed to get details for email ID %s: %v", id, err)
			continue
		}

		// Extract and display relevant information
		fmt.Printf("ID: %s\n", msg.Id)
		for _, header := range msg.Payload.Headers {
			if header.Name == "Subject" || header.Name == "From" {
				fmt.Printf("%s: %s\n", header.Name, header.Value)
			}
		}
		fmt.Printf("Snippet: %s\n\n", msg.Snippet)
	}
}

// BuildQuery constructs a Gmail search query based on various criteria
func BuildQuery(from, to, subject, startDate, endDate string) string {
	var filters []string

	if from != "" {
		filters = append(filters, fmt.Sprintf("from:%s", from))
	}
	if to != "" {
		filters = append(filters, fmt.Sprintf("to:%s", to))
	}
	if subject != "" {
		filters = append(filters, fmt.Sprintf("subject:%s", subject))
	}
	if startDate != "" {
		filters = append(filters, fmt.Sprintf("after:%s", startDate))
	}
	if endDate != "" {
		filters = append(filters, fmt.Sprintf("before:%s", endDate))
	}

	// Join filters with spaces to form the final query string
	return strings.Join(filters, " ")
}

// DisplaySingleEmailDetails fetches and displays details of a specific email by ID
func DisplaySingleEmailDetails(ctx context.Context, gmailTool *tools.GmailTool, emailID string) {
	msg, err := gmailTool.GetEmailDetails(ctx, emailID)
	if err != nil {
		log.Printf("Failed to get details for email ID %s: %v", emailID, err)
		return
	}

	// Extract and display relevant information
	fmt.Printf("ID: %s\n", msg.Id)
	for _, header := range msg.Payload.Headers {
		if header.Name == "Subject" || header.Name == "From" {
			fmt.Printf("%s: %s\n", header.Name, header.Value)
		}
	}
	fmt.Printf("Snippet: %s\n\n", msg.Snippet)
}
