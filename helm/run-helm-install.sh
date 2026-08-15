# #!/bin/bash

# set -e

# NAMESPACE="devops-monitoring-platform"

# kubectl create configmap migration-script \
#   --from-file=init-db.sh=../../docker/init-db.sh \
#   -n "$NAMESPACE"

# kubectl create configmap db-migrations \
#   --from-file=../../backend/migrations \
#   -n "$NAMESPACE"

# helm upgrade --install dev ./helm \
#   -n "$NAMESPACE"
#   -- create-namespace

set -e

NAMESPACE="devops-monitoring-platform"

echo
echo "==========================IMPORTING CONFIGMAP============================"
if kubectl get configmap migration-script -n "$NAMESPACE" >/dev/null 2>&1; then
  echo "ConfigMap migration-script already exists, skipping..."
else
  kubectl create configmap migration-script \
    --from-file=init-db.sh=../docker/init-db.sh \
    -n "$NAMESPACE"
fi

if kubectl get configmap db-migrations -n "$NAMESPACE" >/dev/null 2>&1; then
  echo "ConfigMap db-migrations already exists, skipping..."
else
  kubectl create configmap db-migrations \
    --from-file=../backend/migrations \
    -n "$NAMESPACE"
fi
echo "==========================IMPORTING CONFIGMAP============================"
echo

echo
echo "==========================INSTALL HELM============================"
helm upgrade --install dev ../helm \
  -n "$NAMESPACE" \
  --create-namespace
echo "==========================INSTALL HELM============================"
echo

echo
echo "==========================VERIFY K8S RESOURCES DEPLOYED============================"
kubectl get all -n "$NAMESPACE"
echo "==========================VERIFY K8S RESOURCES DEPLOYED============================"
echo

echo
echo "==========================VERIFY HELM RESOURCES DEPLOYED============================"
helm list --namespace "$NAMESPACE"
echo "==========================VERIFY HELM RESOURCES DEPLOYED============================"
echo