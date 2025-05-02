package filestore_backup

import (
	"context"
	"fmt"
	"log"

	filestore "cloud.google.com/go/filestore/apiv1"
	"cloud.google.com/go/filestore/apiv1/filestorepb"
)

func CreateBackup(ctx context.Context, client *filestore.CloudFilestoreManagerClient, projectID, location, zone, instanceName, backupName, backupDescription string) error {
	instanceFullName := fmt.Sprintf("projects/%s/locations/%s/instances/%s", projectID, zone, instanceName)
	backupFullName := fmt.Sprintf("projects/%s/locations/%s/backups/%s", projectID, zone, backupName)

	// Check if the backup with same name already exists.
	_, err := client.GetBackup(ctx, &filestorepb.GetBackupRequest{Name: backupFullName})
	if err == nil {
		log.Printf("Backup %s already exists. Exiting.\n", backupFullName)
		return nil
	}

	req := &filestorepb.CreateBackupRequest{
		Parent:   fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		BackupId: backupName,
		Backup: &filestorepb.Backup{
			SourceInstance:  instanceFullName,
			Description:     backupDescription,
			SourceFileShare: "airflow",
		},
	}

	fmt.Printf("Creating backup %s from instance %s...\n", backupName, instanceName)
	op, err := client.CreateBackup(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("failed to wait for backup creation: %w", err)
	}

	fmt.Printf("Backup created successfully:\n")
	fmt.Printf("Name: %s\n", resp.Name)
	fmt.Printf("Source Instance: %s\n", resp.SourceInstance)
	fmt.Printf("State: %s\n", resp.State)
	fmt.Printf("Create Time: %s\n", resp.CreateTime)

	return nil
}
