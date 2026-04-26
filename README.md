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
    -v /var/packages/HyperBackup/var/last_result:/var/packages/HyperBackup/var/last_result:ro \
    -v /var/packages/HyperBackup/etc:/var/packages/HyperBackup/etc:ro \
    ghcr.io/ernsthaagsman/hyperbackupexporter:v0.2.1
```

**NOTE** It is necessary to mount the folder, not the individual files. If the files
are mounted, updates to the files will not be seen by the exporter. The files are resolved
to a specific inode when the container is started, and not updated when the files are changed.

Then, in your prometheus configuration, add: 

```yaml
- job_name: 'hyperbackup'
  static_configs:
    - targets:
        - NAShostname:6533
```