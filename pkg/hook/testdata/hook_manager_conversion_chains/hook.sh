#!/usr/bin/env bash

if [[ $1 == "--config" ]] ; then
    cat <<EOF
configVersion: v1
onStartup: 10
kubernetesCustomResourceConversion:
- name: up_conversions
  crdName: crontabs.stable.example.com
  conversions:
  - fromVersion: azaza
    toVersion: ololo
  - fromVersion: ololo
    toVersion: foobar
  - fromVersion: foobar
    toVersion: abc
- name: down_conversions
  crdName: crontabs.stable.example.com
  conversions:
  - fromVersion: abc
    toVersion: foobar
  - fromVersion: foobar
    toVersion: ololo
  - fromVersion: ololo
    toVersion: azaza
EOF

else
  exit 0
fi