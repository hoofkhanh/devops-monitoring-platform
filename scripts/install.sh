#!/bin/bash
set -euo pipefail

sudo -v || exit 1

# Set up a systemd service to run the monitoring agent script in the background
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

sudo tee /etc/systemd/system/monitoring-agent.service > /dev/null <<EOF
[Unit]
Description=Devops Monitoring Agent
After=network.target

[Service]
Type=simple
ExecStart=$PROJECT_DIR/monitoring-agent/collect-metrics.sh
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

sudo chmod u+x "$PROJECT_DIR/monitoring-agent/collect-metrics.sh"
sudo systemctl daemon-reload
sudo systemctl start monitoring-agent

# Run helm install script to deploy the application
sudo chmod u+x "$PROJECT_DIR/helm/run-helm-install.sh"
cd "$PROJECT_DIR/helm"
./run-helm-install.sh