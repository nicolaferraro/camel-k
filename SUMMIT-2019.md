# Summit 2019 Instructions

This repo provides the version of Camel K used for the Summit demo.

## Base Syndesis image

The `./base-syndesis/release.sh` script publishes a new version of the base image for Syndesis.

That will not be required on Knative Serving 0.4+.

## Creating Camel K images

**NOTE**: This is the last image that should be build. Camel K Runtime SNAPSHOTS and Syndesis SNAPSHOTS
needs to be installed in the local maven repo before building the Camel K image

```
# Building the image
make clean images-dev

# Pushing to quay.io (you need to join the redhatdemo team)
docker push quay.io/redhatdemo/camel-k:latest
```
