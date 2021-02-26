# ads-service-module

`ads-service-module` is a Azure IoT Hub module that connects and sends events to the Azure IoT Hub.

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
