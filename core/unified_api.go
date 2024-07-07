package core

type CloudService interface {
	UploadFile(filePath string) error
	DownloadFile(filePath string) error
	// Add more methods as needed
}
