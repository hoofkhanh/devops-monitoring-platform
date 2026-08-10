kubectl create configmap migration-script \
  --from-file=init-db.sh=../../docker/init-db.sh \
  -n devops-monitoring-platform

kubectl create configmap db-migrations \
  --from-file=../../backend/migrations \
  -n devops-monitoring-platform