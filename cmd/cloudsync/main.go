package main

import (
	"github.com/mczimm/cloud-synchronizer/services/icloud"
	"log"
	"os"

	"github.com/mczimm/cloud-synchronizer/core"
	"github.com/mczimm/cloud-synchronizer/services/google_drive"
)

func main() {
	// Configure loggers to include file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) < 3 {
		log.Println("usage: cloudsync <from> <to>")
		return
	}

	fromCloud := os.Args[2]
	toCloud := os.Args[3]

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
