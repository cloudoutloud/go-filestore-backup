package main

import (
	"context"
	"filestore-backup/filestore_backup"
	"log"
	"os"
	"time"

	filestore "cloud.google.com/go/filestore/apiv1"
)

var (
	projectID    = os.Getenv("GCP_PROJECT_ID")
	location     = os.Getenv("GCP_LOCATION")
	zone         = os.Getenv("GCP_ZONE")
	instanceName = os.Getenv("GCP_FILESTORE_INSTANCE")
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// No specific auth set will look for GOOGLE_APPLICATION_CREDENTIALS
	client, err := filestore.NewCloudFilestoreManagerClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	backupName := "automated-backup-" + time.Now().Format("2006-01-02")
	backupDescription := "Automated daily backup"

	if err := filestore_backup.CreateBackup(ctx, client, projectID, location, zone, instanceName, backupName, backupDescription); err != nil {
		log.Printf("Failed to create backup: %v", err)
	}

	// 7 days
	olderThan := time.Duration(168) * time.Hour

	if err := filestore_backup.DeleteOldBackups(ctx, client, projectID, location, olderThan); err != nil {
		log.Fatalf("Failed to delete old backups: %v", err)
	}
}
