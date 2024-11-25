// Todos:
// Please implement a gmail tool that
// ran list (list of email Ids). done
// read email. done
// move emails.
// Bonus would be list email with filter (e.g. from, to, subject, time range etc) in progress
package main

import (
	"bufio"
	"context"
	"fmt"
	"gmail_tool_project/helpers"
	"gmail_tool_project/tools"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	ctx := context.Background()

	fmt.Println("Select from the following options:")
	fmt.Println("1. Create-Token")
	fmt.Println("2. List EmailIds in Inbox")
	fmt.Println("3. Enter EmailID to Read Email")
	fmt.Println("4. List Emails with Filter")

	// Read user input
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter option: ")
	input, _ := reader.ReadString('\n')

	// Convert input to an integer
	option, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		log.Fatalf("Invalid option: %v", err)
	}

	switch option {
	case 1:
		// Create token by running the OAuth flow
		err := tools.CreateToken(ctx)
		if err != nil {
			log.Fatalf("Failed to create token: %v", err)
		}
		fmt.Println("Token created successfully and saved to token.json")

	case 2:
		// List email IDs in inbox
		gmailTool := helpers.InitializeGmailTool(ctx)
		helpers.DisplayEmailIDs(ctx, gmailTool, "in:inbox")

	case 3:
		// List email IDs in inbox
		gmailTool := helpers.InitializeGmailTool(ctx)
		helpers.DisplayEmailIDs(ctx, gmailTool, "in:inbox")

		// Prompt the user to enter an email ID to read
		fmt.Print("Enter the Email ID to read: ")
		emailID, _ := reader.ReadString('\n')
		emailID = strings.TrimSpace(emailID)

		// Display email details for the provided ID
		helpers.DisplaySingleEmailDetails(ctx, gmailTool, emailID)

	case 4:
		// List emails with filters
		gmailTool := helpers.InitializeGmailTool(ctx)

		// Get filtering criteria from user input
		fmt.Println("Enter filter criteria. Leave blank to skip a filter.")
		fmt.Print("From (email address): ")
		from, _ := reader.ReadString('\n')
		fmt.Print("To (email address): ")
		to, _ := reader.ReadString('\n')
		fmt.Print("Subject: ")
		subject, _ := reader.ReadString('\n')
		fmt.Print("Start Date (YYYY/MM/DD): ")
		startDate, _ := reader.ReadString('\n')
		fmt.Print("End Date (YYYY/MM/DD): ")
		endDate, _ := reader.ReadString('\n')

		// Trim whitespace and prepare query string
		query := helpers.BuildQuery(
			strings.TrimSpace(from),
			strings.TrimSpace(to),
			strings.TrimSpace(subject),
			strings.TrimSpace(startDate),
			strings.TrimSpace(endDate),
		)

		helpers.DisplayEmailDetails(ctx, gmailTool, query)

	default:
		fmt.Println("Invalid option selected.")
	}
}
