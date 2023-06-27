#!/bin/bash
set -xeou pipefail

export TAG=$1

yq e -i '.adminApi.image.tag = strenv(TAG)' helm/kai/values.yaml
yq e -i '.k8sManager.image.tag = strenv(TAG)' helm/kai/values.yaml
yq e -i '.mongoWriter.image.tag = strenv(TAG)' helm/kai/values.yaml
yq e -i '.natsManager.image.tag = strenv(TAG)' helm/kai/values.yaml
yq e -i '.version = strenv(TAG)' helm/kai/Chart.yaml
yq e -i '.appVersion = strenv(TAG)' helm/kai/Chart.yaml
