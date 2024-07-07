package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mczimm/cloud-synchronizer/services/icloud"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"

	"github.com/mczimm/cloud-synchronizer/core"
	"github.com/mczimm/cloud-synchronizer/services/google_drive"

	"golang.org/x/oauth2/google"
)

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("unable to retrieve token from web %v", err)
	}
	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("can't close file %v", err)
		}
	}(f)
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("unable to cache oauth token: %v", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("can't close file %v", err)
		}
	}(f)
	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		log.Fatalf("unable to encode token %v", err)
	}
}

func main() {
	// Configure loggers to include file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) < 3 {
		log.Println("usage: cloudsync <creds> <from> <to>")
		return
	}

	creds := os.Args[2]
	fromCloud := os.Args[3]
	toCloud := os.Args[4]

	ctx := context.Background()
	b, err := os.ReadFile(creds)
	if err != nil {
		log.Printf("unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
	if err != nil {
		log.Fatalf("unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	_, err = drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("unable to retrieve Drive client: %v", err)
	}

	syncManagerFromCloud := core.NewSyncManager()
	syncManagerFromCloud.RegisterService("google_drive", new(google_drive.GoogleDriveAdapter))

	_, exists := syncManagerFromCloud.GetService(fromCloud)
	if !exists {
		log.Printf("cloud service '%s' not registered", fromCloud)
		os.Exit(1)
	}

	syncManagerToCloud := core.NewSyncManager()
	syncManagerToCloud.RegisterService("icloud", new(icloud.GoogleDriveAdapter))

	_, exists = syncManagerToCloud.GetService(toCloud)
	if !exists {
		log.Printf("cloud service '%s' not registered", toCloud)
	}
}
