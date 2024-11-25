package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type GmailTool struct {
	client *gmail.Service
}

// InitializeGmailTool initializes GmailTool with an authenticated client
func InitializeGmailTool(ctx context.Context) (*GmailTool, error) {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %v", err)
	}

	// Configure OAuth2 with credentials
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)

	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Gmail client: %v", err)
	}

	return &GmailTool{client: service}, nil
}

// ListEmailIDs fetches email IDs from Gmail based on a query
func (g *GmailTool) ListEmailIDs(ctx context.Context, query string) ([]string, error) {
	user := "me"
	var emailIDs []string

	req := g.client.Users.Messages.List(user).Q(query)
	err := req.Pages(ctx, func(page *gmail.ListMessagesResponse) error {
		for _, msg := range page.Messages {
			emailIDs = append(emailIDs, msg.Id)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list email IDs: %w", err)
	}

	return emailIDs, nil
}

// Supporting function to handle OAuth2
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser: \n%v\n", authURL)

	codeCh := make(chan string)

	// Start a local web server to handle the callback
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		codeCh <- code
		fmt.Fprintln(w, "Authorization successful! You may close this window.")
	})

	server := &http.Server{Addr: "localhost:8080"}
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	code := <-codeCh
	if err := server.Close(); err != nil {
		log.Printf("Error closing server: %v", err)
	}

	// Exchange code for token
	tok, err := config.Exchange(context.TODO(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Function to retrieve a token from a local file
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Function to save a token to a file
func saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to cache oauth token: %v", err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}

// GetEmailDetails fetches basic details of an email by its ID
func (g *GmailTool) GetEmailDetails(ctx context.Context, emailID string) (*gmail.Message, error) {
	msg, err := g.client.Users.Messages.Get("me", emailID).Format("full").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve email details for ID %s: %w", emailID, err)
	}
	return msg, nil
}

// CreateToken initiates the OAuth flow to create and save a new token
func CreateToken(ctx context.Context) error {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return fmt.Errorf("unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		return fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	// Run the OAuth flow and save the token
	tok := getTokenFromWeb(config)
	return saveToken("token.json", tok)
}
