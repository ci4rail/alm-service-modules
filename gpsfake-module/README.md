# gpsfake-module

`gpsfake-module` is a module that provides simulated location data provided by a `gpsd` server.
The `gpsfake-module` behaves like gpsd server simulating a gps modem. Use the `gpsfake-module` with `alm-location-module`.

*Note: This is based on https://github.com/knowhowlab/gpsd-nmea-simulator*

## Example usage
Create a new deployment manifest that contains the module, e.g. `myapplication.yaml`.

**Note: when connecting to another module e.g. from `location` to `gpsfake` the hostname to connect to is in the form of:   `<application-name>_<module-name>`, e.g. `my-application_gpsfake`.**

```yaml
---
application: my-application
modules:
  - name: location
    image: ci4rail/alm-location-module:latest
    createOptions: '{}'
    imagePullPolicy: on-create
    restartPolicy: always
    status: running
    startupOrder: 1
    env:
      GPSD_HOST=my-application_gpsfake
  - name: gpsfake
    image: ci4rail/gpsfake-module:latest
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
