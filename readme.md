# ScalewayBackup

ScalewayBackup is a small application automating process of making backup copies
of [Scaleway](https://www.scaleway.com) machines using snapshots.

ScalewayBackup has been designed as a first-class Kubernetes citizen.

## Usage

The easiest way to install ScalewayBackup is to execute Helm chart provided with the project:

    git clone git@github.com:hekonsek/scaleway-backup.git
    cd scaleway-backup
    helm upgrade scaleway-backup helm --install \
     --set token=myScalewayToken --set organization=myScalewayOrganizationId \
     --set 'volumes=volume1_id\,volume2_id'