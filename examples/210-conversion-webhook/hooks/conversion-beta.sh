#!/usr/bin/env bash

source /shell_lib.sh

function __config__(){
    cat <<EOF
configVersion: v1
kubernetesCustomResourceConversion:
  - name: crontabs
    crdName: crontabs.stable.example.com
    conversions:
    - fromVersion: v1beta1
      toVersion: v1beta2
    - fromVersion: v1beta2
      toVersion: v1
EOF
}

function __crd_conversion::crontabs::v1beta1_to_v2beta2() {
  image=$(context::jq -r '.review.request.object.spec.image')
  echo "Got image: $image"

  if [[ $image == repo.example.com* ]] ; then
    cat <<EOF > $VALIDATING_RESPONSE_PATH
{"allowed":true}
EOF
  else
    cat <<EOF > $VALIDATING_RESPONSE_PATH
{"allowed":false, "message":"Only images from repo.example.com are allowed"}
EOF
  fi
}

function __crd_conversion::crontabs::v1beta2_to_v1() {
  image=$(context::jq -r '.review.request.object.spec.image')
  echo "Got image: $image"

  if [[ $image == repo.example.com* ]] ; then
    cat <<EOF > $VALIDATING_RESPONSE_PATH
{"allowed":true}
EOF
  else
    cat <<EOF > $VALIDATING_RESPONSE_PATH
{"allowed":false, "message":"Only images from repo.example.com are allowed"}
EOF
  fi
}

function __main__() {
      cat <<EOF >$CONVERSION_RESPONSE_PATH
{"failedMessage":"Conversion for CRONTABS beta is not implemented yet"}
EOF

}

hook::run $@
