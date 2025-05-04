# go-filestore-backup

Script to automate backups of GCP-managed NFS [filestore](https://cloud.google.com/filestore)

Although backups are supported, these are not automated out of the box like managed DB offerings.

This script will create a backup of a filestore instance named `automated-backup-<date-run>`
Then will look and delete backups older than the set threshold to reduce cost.

## Example build/push image

There is no automation currently around building and pushing the image to a private registry.

Manual example using Docker cli, from the root of the repo.

`docker build . -t filestore-backup:v1.0.1`

Then using `docker push` to push to the registry.

## Deployment

You normally only have access to the filestore from within a set VPC network.
In the `k8s` folder then is example manifests on how to deploy this as a cron job and run within a Kubernetes clusters residing in the same VPC as the filestore.

Alternatively, you could run this in serverless offerings such as:
[cloudrun](https://cloud.google.com/run)

## Authentication to GCP

Script is using [ADC](https://cloud.google.com/docs/authentication/provide-credentials-adc) to authenticate.

You can create a [GCP service account](https://cloud.google.com/iam/docs/service-account-overview) and grant the role [file.editor](https://cloud.google.com/iam/docs/understanding-roles#cloud-filestore-roles)

You can then proceed to generate a service account JSON.

## Environment vars needed

`GCP_PROJECT_ID`- GCP project ID where the filestore instance is located

`GCP_LOCATION` - Region name where the filestore instance is located

`GCP_ZONE` - Zone name where the filestore instance is located

`GCP_FILESTORE_INSTANCE` - The name of the filestore instance

`GCP_FILESTORE_SHARE_NAME` - The Name of the share on the filestore instance

`BACKUP_DURATION` - Set threshold to deleted backups older than this value in hours (default 168 days/ 7 days)

`GOOGLE_APPLICATION_CREDENTIALS` - The path to the GCP service account JSON key used for authentication (Not ideal for now as it uses long-lived keys)

