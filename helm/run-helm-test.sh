#!/bin/bash

NAMESPACE="devops-monitoring-platform"

set -e

read -rp "Enter backend hostname: " HOSTNAME
HOSTNAME=${HOSTNAME:-hoofkhanh.com}

echo
echo "==========================HELM LINT============================"
helm lint ../helm
echo "==========================HELM LINT============================"
echo

echo
echo "==========================HELM TEMPLATE============================"
helm template ../helm \
    -n "$NAMESPACE" \
	--set ingress.host="$HOSTNAME"
echo "==========================HELM TEMPLATE============================"
echo

echo
echo "==========================IMPORTING CONFIGMAP IN A SIMULATION ENV============================"
if kubectl get configmap migration-script -n "$NAMESPACE" >/dev/null 2>&1; then
	echo "ConfigMap migration-script already exists, skipping..."
else
	kubectl create configmap migration-script \
  	--from-file=init-db.sh=../docker/init-db.sh \
  	-n "$NAMESPACE" \
  	--dry-run=server
fi

if kubectl get configmap db-migrations -n "$NAMESPACE" >/dev/null 2>&1; then
	echo "ConfigMap db-migrations already exists, skipping..."
else
	kubectl create configmap db-migrations \
  	--from-file=../backend/migrations \
  	-n "$NAMESPACE" \
  	--dry-run=server
fi



echo "==========================IMPORTING CONFIGMAP IN A SIMULATION ENV============================"
echo

echo
echo "==========================INSTALL HELM============================"
helm upgrade --install dev ../helm \
  -n "$NAMESPACE" \
  --create-namespace \
  --set ingress.host="$HOSTNAME" \
  --dry-run=server
echo "==========================INSTALL HELM============================"
echo