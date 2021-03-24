# alm-location-module

`alm-location-module` is a module provides location data by connecting to a `gpsd` server.

## Example usage
Create a new deployment manifest that contains the module, e.g. `myapplication.yaml`.

```yaml
---
application: my-application
modules:
  - name: alm-location-module
    image: alm-location-module:latest
    createOptions: '{}'
    imagePullPolicy: on-create
    restartPolicy: always
    status: running
    startupOrder: 1
```

Deploy the manifest using the kyt-cli

```sh
$ kyt alm apply -f myapplication.yaml
```
