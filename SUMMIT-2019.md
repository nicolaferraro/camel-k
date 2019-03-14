# Summit 2019 Instructions

This repo provides the version of Camel K used for the Summit demo.

## Creating the images

**NOTE**: This is the last image that should be build. Camel K Runtime SNAPSHOTS and Syndesis SNAPSHOTS
needs to be installed in the local maven repo before building the Camel K image

```
# Building the image
make clean images-dev

# Pushing to quay.io (you need to join the redhatdemo team)
docker push quay.io/redhatdemo/camel-k:latest
```
