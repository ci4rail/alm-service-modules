# alm-location-module

`alm-location-module` is a module provides location data by connecting to a `gpsd` server.

## Example usage
Create a new deployment manifest that contains the module, e.g. [`manifest.yaml`](example/manifest.yaml).
**Note: replace the `<TAG>` field with a valid tag for the image.**

Deploy the manifest using the edgefarm-cli

```sh
$ edgefarm alm apply -f example/manifest.yaml
```
