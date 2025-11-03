#!/bin/bash
BACKUP_DIR="/mnt/k8s-storage/backups"
DATE=$(date +%Y%m%d_%H%M%S)

echo "🔄 Iniciando backup..."

# Backup PostgreSQL
echo "📦 Haciendo backup de PostgreSQL..."
kubectl exec deployment/postgres-deployment -- pg_dump -U postgres myapp > $BACKUP_DIR/postgres_$DATE.sql

# Backup Redis
echo "📦 Haciendo backup de Redis..."
kubectl exec deployment/redis-deployment -- redis-cli SAVE
kubectl cp $(kubectl get pods -l app=redis -o name | cut -d/ -f2):/data/dump.rdb $BACKUP_DIR/redis_$DATE.rdb

# Comprimir
echo "🗜️ Comprimiendo backups..."
tar -czf $BACKUP_DIR/backup_$DATE.tar.gz -C $BACKUP_DIR postgres_$DATE.sql redis_$DATE.rdb

# Limpiar archivos temporales
rm $BACKUP_DIR/postgres_$DATE.sql $BACKUP_DIR/redis_$DATE.rdb

echo "✅ Backup completado: $BACKUP_DIR/backup_$DATE.tar.gz"