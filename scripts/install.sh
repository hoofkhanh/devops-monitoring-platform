
sudo nano /etc/systemd/system/monitoring-agent.service
[Unit]
Description=Devops Monitoring Agent
After=network.target

[Service]
Type=simple
ExecStart=/home/dell/pet-projecs/devops-monitoring-platform/monitoring-agent/collect-metrics.sh
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target


sudo chmod u+x  /etc/systemd/system/monitoring-agent.service

sudo systemctl daemon-reload
sudo systemctl enable monitoring-agent
sudo systemctl start monitoring-agent