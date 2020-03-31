#!/bin/sh

location=$(dirname $0)

PACKAGE=${1:-camel-k-nightly}
VERSION=${2:-1.0.0-nightly.202003261646}
ORGANIZATION=${3:-nferraro}

echo "This script can be used to upload the operator manifests to Quay.io"
echo "Before publishing:"
echo "- Set the environment variables QUAY_USERNAME and QUAY_PASSWORD (e.g. using envrc)"
echo "- If already present, the operator should be removed from the \"Applications\" tab on Quay.io"
echo ""
echo "After publication:"
echo "- Make the application public on Quay.io (it's private by default)"
echo ""
echo "To use the operator in a cluster:"
echo "- Ensure the integration operators source is installed in the cluster (oc apply -f operator-source.yaml)"
echo ""
echo ""
echo "Now publishing the $PACKAGE manifests (version $VERSION) to quay.io/$ORGANIZATION ..."
echo ""

if [ -z "$QUAY_USERNAME" ]; then
  echo "QUAY_USERNAME environment variable cannot be empty"
  exit 1
fi

if [ -z "$QUAY_PASSWORD" ]; then
  echo "QUAY_PASSWORD environment variable cannot be empty"
  exit 1
fi

export AUTH_TOKEN=$(curl -sH "Content-Type: application/json" -XPOST https://quay.io/cnr/api/v1/users/login -d '{"user": {"username": "'"${QUAY_USERNAME}"'", "password": "'"${QUAY_PASSWORD}"'"}}' | jq -r '.token')

cd $location
operator-courier --verbose push ./camel-k-nightly ${ORGANIZATION} ${PACKAGE} ${VERSION} "$AUTH_TOKEN"
