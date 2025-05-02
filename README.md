# go-filestore-backup

Script to automate backups of google filestore

Will create a backup of a filestore instance named `automated-backup-<date-run>`

Will look and delete backups older than default 7 days

## Environment vars needed

`GCP_PROJECT_ID`- GCP project id where filestore instance is located
`GCP_LOCATION` - region name where filestore instance is located
`GCP_ZONE` - zone name where filestore instance is located
`GCP_FILESTORE_INSTANCE` - name of filestore instance
`GCP_FILESTORE_SHARE_NAME` - name of share on filestore instance
`BACKUP_DURATION` - deleted backups older than this duration
`GOOGLE_APPLICATION_CREDENTIALS` - path to authentication google service account JSON key

