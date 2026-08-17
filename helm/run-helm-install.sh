#!/bin/bash

set -e

read -rp "Enter backend hostname: " HOSTNAME
HOSTNAME=${HOSTNAME:-hoofkhanh.com}

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
  --set ingress.host="$HOSTNAME" \
  --create-namespace


kubectl rollout restart deployment "dev-${NAMESPACE}-frontend" -n "$NAMESPACE"
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

# config the ingress host into the /etc/hosts file
INGRESS_ADDRESS=$(kubectl get ingress dev-devops-monitoring-platform-ingress \
  -n devops-monitoring-platform -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
ESCAPED_IP=$(printf '%s\n' "$INGRESS_ADDRESS" | sed 's/\./\\./g')

if grep -q "^${ESCAPED_IP}" /etc/hosts; then
    sudo sed -i "s/^${ESCAPED_IP}.*$/${INGRESS_ADDRESS} $HOSTNAME/" /etc/hosts
else
    sudo sh -c "echo '$INGRESS_ADDRESS $HOSTNAME' >> /etc/hosts"
fi

sed -i "s|http://[^/]*|http://$HOSTNAME:80|" ../monitoring-agent/collect-metrics.sh

echo "url: http://$HOSTNAME"