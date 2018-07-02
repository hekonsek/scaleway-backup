# ScalewayBackup

ScalewayBackup is a small application automating process of making backup copies
of [Scaleway](https://www.scaleway.com) machines using snapshots.

ScalewayBackup has been designed as a first-class Kubernetes citizen.

## How it works?

ScalewayBackup is a small application written in GoLang. When you execute the app, it creates the snapshot of volumes
indicated to be backed-up. By default Helm chart deploys ScalewayBackup as a cron job executed once a day (in the middle
of the night). It holds 10 last copies of the backup for every volume (older backup copies are deleted). So effectively
default setup holds backup copies of your volumes from last ten days.

## Usage

The easiest way to install ScalewayBackup is to execute Helm chart provided with the project:

    git clone git@github.com:hekonsek/scaleway-backup.git
    helm upgrade scaleway-backup scaleway-backup/charts/scaleway-backup --install \
     --set token=myScalewayToken --set organization=myScalewayOrganizationId \
     --set 'volumes=volume1_id\,volume2_id'