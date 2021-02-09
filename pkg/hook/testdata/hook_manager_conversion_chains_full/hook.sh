#!/usr/bin/env bash

if [[ $1 == "--config" ]] ; then
    cat <<EOF
configVersion: v1
onStartup: 10
kubernetesCustomResourceConversion:
- name: up_conversions
  crdName: crontabs.stable.example.com
  conversions:
  - fromVersion: group.io/azaza
    toVersion: ololo
  - fromVersion: unstable.example.com/ololo
    toVersion: foobar
  - fromVersion: stable.example.com/foobar
    toVersion: next.io/abc
- name: down_conversions
  crdName: crontabs.stable.example.com
  conversions:
  - fromVersion: stable.example.com/abc
    toVersion: stable.example.com/foobar
  - fromVersion: foobar
    toVersion: ololo
  - fromVersion: ololo
    toVersion: azaza
EOF

else
  exit 0
fi