HyperBackup Prometheus Exporter
===============================

This repo contains a basic Prometheus exporter for Synology HyperBackup.

As of now, all it exposes is the date of the last backup, and the incremental version of it.

If you have further ideas for expanding this exporter, feel free to open an issue or a PR.

How to run
----------

On the Synology NAS, with Docker installed, connect using SSH and run:

```
docker run -d \
    -p 6533:6533 \
    --name hyperbackup-exporter \
    -v /var/packages/HyperBackup/var/last_result/backup.last:/var/packages/HyperBackup/var/last_result/backup.last:ro \
    -v /var/packages/HyperBackup/etc/synobackup.conf:/var/packages/HyperBackup/etc/synobackup.conf:ro \
    ghcr.io/ernsthaagsman/hyperbackupexporter:v0.1
```

Then, in your prometheus configuration, add: 

```yaml
- job_name: 'hyperbackup'
  static_configs:
    - targets:
        - NAShostname:6533
```