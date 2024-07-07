package core

type CloudService interface {
	UploadFile(filePath string) error
	DownloadFile(filePath string) error
	SyncFolder(fromFolder, toFolder string) error
	// Add more methods as needed
}
