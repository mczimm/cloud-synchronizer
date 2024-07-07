package main

import (
	"fmt"
	"github.com/mczimm/cloud-synchronizer/core"
	"github.com/mczimm/cloud-synchronizer/services/google_drive"
	"os"
	// Import other cloud service adapters as needed
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: uploadfile <cloud_service_name> <file_path>")
		os.Exit(1)
	}

	cloudServiceName := os.Args[1]
	filePath := os.Args[2]

	syncManager := core.NewSyncManager()

	// Register cloud service adapters
	syncManager.RegisterService("google_drive", new(google_drive.GoogleDriveAdapter))
	// Register other services as needed

	service, exists := syncManager.GetService(cloudServiceName)
	if !exists {
		fmt.Printf("Cloud service '%s' not registered\n", cloudServiceName)
		os.Exit(1)
	}

	err := service.UploadFile(filePath)
	if err != nil {
		fmt.Printf("Failed to upload file: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("File uploaded successfully")
}
