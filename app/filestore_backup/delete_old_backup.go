package filestore_backup

import (
	"context"
	"fmt"
	"log"
	"time"

	filestore "cloud.google.com/go/filestore/apiv1"
	"cloud.google.com/go/filestore/apiv1/filestorepb"
	"google.golang.org/api/iterator"
)

func DeleteOldBackups(ctx context.Context, client *filestore.CloudFilestoreManagerClient, projectID, location string, olderThan time.Duration) error {
	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, location)

	req := &filestorepb.ListBackupsRequest{
		Parent: parent,
	}

	it := client.ListBackups(ctx, req)
	for {
		backup, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to list backups: %w", err)
		}

		// Check if the backup is older than the specified duration
		backupAge := time.Since(backup.CreateTime.AsTime())
		if backupAge > olderThan {
			backupName := backup.Name
			fmt.Printf("Deleting backup %s (age: %s)...\n", backupName, backupAge)

			deleteReq := &filestorepb.DeleteBackupRequest{
				Name: backupName,
			}

			op, err := client.DeleteBackup(ctx, deleteReq)
			if err != nil {
				log.Printf("Failed to delete backup %s: %v\n", backupName, err)
				continue
			}

			err = op.Wait(ctx)
			if err != nil {
				log.Printf("Failed to wait for deletion of backup %s: %v\n", backupName, err)
				continue
			}

			fmt.Printf("Backup %s deleted successfully.\n", backupName)
		} else {
			fmt.Printf("Skipping backup %s (age: %s), not older than %s.\n", backup.Name, backupAge, olderThan)
		}
	}

	return nil
}
