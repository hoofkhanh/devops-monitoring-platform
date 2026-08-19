#!/bin/bash

echo "==============================================================================="
echo "Monitoring Agent - Collecting Metrics"
echo "==============================================================================="
echo

set -euo pipefail

while true; do
    echo "==============================================================================="
    cpu_usage=$(top -bn 1 | grep "%Cpu" | awk '{printf "%.2f\n", 100 - $8}')
    echo "CPU USAGE: $cpu_usage%"

    memory_usage=$(free -h | grep "Mem:" | awk '{printf "%.2f\n", ($3 / $2) * 100}')
    echo "MEMORY USAGE: $memory_usage%"

    disk_usage=$(df -h | grep "/dev/sdd" | awk '{print $5}' | sed 's/%//')
    echo "DISK USAGE: $disk_usage%"
    echo
    
    # 80 is a nginx's port, which is a reverse proxy for the backend API
    curl -X POST http://khanh.com:80/api/metrics \
        -H "Content-type: application/json" \
        -d '{"cpu": '"$cpu_usage"', "memory": '"$memory_usage"', "disk": '"$disk_usage"'}'

    echo "Metrics sent to backend API."
    echo "==============================================================================="
    echo

    sleep 5
done