#!/bin/bash
set -euo pipefail

sudo -v || exit 1

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

INGRESS_IP=$(kubectl get ingress dev-devops-monitoring-platform-ingress \
    -n devops-monitoring-platform -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
ESCAPED_IP=$(printf '%s\n' "$INGRESS_IP" | sed 's/\./\\./g')    

sudo sed -i "/^${ESCAPED_IP}.*$/d" /etc/hosts

sudo systemctl stop monitoring-agent 
sudo systemctl daemon-reload

echo
echo "==========================VERIFY MONITORING AGENT============================"
sudo systemctl status monitoring-agent || true
echo "==========================VERIFY MONITORING AGENT============================"
echo

sudo chmod u+x "$PROJECT_DIR/helm/run-helm-uninstall.sh"
cd "$PROJECT_DIR/helm"
./run-helm-uninstall.sh