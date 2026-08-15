#!/bin/bash

NAMESPACE="devops-monitoring-platform"

set -e

echo
echo "==========================HELM LINT============================"
helm lint ../helm
echo "==========================HELM LINT============================"
echo

echo
echo "==========================HELM TEMPLATE============================"
helm template ../helm \
    -n "$NAMESPACE" --debug
echo "==========================HELM TEMPLATE============================"
echo

echo
echo "==========================IMPORTING CONFIGMAP IN A SIMULATION ENV============================"
kubectl create configmap migration-script \
  --from-file=init-db.sh=../docker/init-db.sh \
  -n "$NAMESPACE" \
  --dry-run=server -o yaml

kubectl create configmap db-migrations \
  --from-file=../backend/migrations \
  -n "$NAMESPACE" \
  --dry-run=server -o yaml
echo "==========================IMPORTING CONFIGMAP IN A SIMULATION ENV============================"
echo

echo
echo "==========================INSTALL HELM============================"
helm upgrade --install dev ../helm \
  -n "$NAMESPACE" \
  --create-namespace \
  --dry-run=server
echo "==========================INSTALL HELM============================"
echo