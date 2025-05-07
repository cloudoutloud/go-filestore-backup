package main

import (
	"context"
	"filestore-backup/filestore_backup"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	filestore "cloud.google.com/go/filestore/apiv1"
)

var (
	projectID      = os.Getenv("GCP_PROJECT_ID")
	region         = os.Getenv("GCP_REGION")
	zone           = os.Getenv("GCP_ZONE")
	instanceName   = os.Getenv("GCP_FILESTORE_INSTANCE")
	fileshareName  = os.Getenv("GCP_FILESTORE_SHARE_NAME")
	backupDuration = os.Getenv("BACKUP_DURATION")
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

	if err := filestore_backup.CreateBackup(ctx, client, projectID, region, zone, instanceName, backupName, backupDescription, fileshareName); err != nil {
		log.Printf("Failed to create backup: %v", err)
	}

	if backupDuration == "" {
		fmt.Println("DURATION environment variable not set.  Using default of 168 hours (7 days).")
		backupDuration = "168"
	}

	durationHours, err := strconv.Atoi(backupDuration)
	if err != nil {
		fmt.Printf("Error converting DURATION to integer: %v.  Using default of 168 hours (7 days).\n", err)
		durationHours = 168
	}

	olderThan := time.Duration(durationHours) * time.Hour

	if err := filestore_backup.DeleteOldBackups(ctx, client, projectID, region, olderThan); err != nil {
		log.Fatalf("Failed to delete old backups: %v", err)
	}
}
