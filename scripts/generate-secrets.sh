#!/bin/bash

# scripts/generate-secrets.sh
set -e  # Exit on error

echo "🔐 Generating secure secrets from your .env file..."

# Cargar variables del .env si existe
if [ -f .env ]; then
    echo "📝 Loading environment from .env file..."
    export $(cat .env | grep -v '^#' | xargs)
else
    echo "⚠️  No .env file found, using defaults or generating random values"
fi

# Función para obtener valor con compatibilidad de nombres
get_secure_value() {
    local primary_var=$1
    local fallback_var=$2
    local default_length=$3
    local description=$4
    
    local value=""
    
    # Verificar variable primaria primero, luego fallback
    if [ -n "${!primary_var}" ]; then
        value="${!primary_var}"
    elif [ -n "${!fallback_var}" ]; then
        value="${!fallback_var}"
    else
        value=$(openssl rand -base64 "$default_length" | head -c "$default_length")
    fi
    
    echo "$value"
}

# Función para obtener valor con SHA256
get_secure_value_with_hash() {
    local primary_var=$1
    local fallback_var=$2
    local default_length=$3
    local description=$4
    
    local value=""
    
    # Verificar variable primaria primero, luego fallback
    if [ -n "${!primary_var}" ]; then
        echo "   ✅ Using $primary_var: $description"
        value="${!primary_var}"
    elif [ -n "${!fallback_var}" ]; then
        echo "   ✅ Using $fallback_var (fallback): $description"
        value="${!fallback_var}"
    else
        echo "   ⚠️  Neither $primary_var nor $fallback_var set, generating random value"
        value=$(openssl rand -base64 "$default_length" | head -c "$default_length")
    fi
    
    # Generar SHA256 del valor
    local hash
    hash=$(echo -n "$value" | sha256sum | cut -d' ' -f1)
    
    # Devolver valor y hash
    echo "$value|$hash"
}

echo ""
echo "📝 Generating values from your configuration..."

# Obtener valores con compatibilidad de nombres
JWT_DATA=$(get_secure_value_with_hash "JWT_SECRETO" "JWT_SECRET" 32 "JWT Secret")
JWT_SECRET=$(echo "$JWT_DATA" | cut -d'|' -f1)
JWT_SECRET_HASH=$(echo "$JWT_DATA" | cut -d'|' -f2)

DB_DATA=$(get_secure_value_with_hash "PASS_DB" "DB_PASSWORD" 24 "Database Password")
DB_PASSWORD=$(echo "$DB_DATA" | cut -d'|' -f1)
DB_PASSWORD_HASH=$(echo "$DB_DATA" | cut -d'|' -f2)

REDIS_DATA=$(get_secure_value_with_hash "REDIS_PASSWORD" "REDIS_PASS" 16 "Redis Password")
REDIS_PASSWORD=$(echo "$REDIS_DATA" | cut -d'|' -f1)
REDIS_PASSWORD_HASH=$(echo "$REDIS_DATA" | cut -d'|' -f2)

JWT_SALT_DATA=$(get_secure_value_with_hash "JWT_SALT" "JWT_SALT_ROUNDS" 32 "JWT Salt")
JWT_SALT=$(echo "$JWT_SALT_DATA" | cut -d'|' -f1)
JWT_SALT_HASH=$(echo "$JWT_SALT_DATA" | cut -d'|' -f2)

# Obtener otras configuraciones importantes
DB_USER=$(get_secure_value "USER_DB" "DB_USER" 0 "Database User")
DB_NAME=$(get_secure_value "NAME_DB" "DB_NAME" 0 "Database Name")
DB_HOST=$(get_secure_value "HOST_DB" "DB_HOST" 0 "Database Host")
DB_PORT=$(get_secure_value "PORT_DB" "DB_PORT" 0 "Database Port")
REDIS_HOST=$(get_secure_value "REDIS_URL" "REDIS_HOST" 0 "Redis Host")
REDIS_PORT=$(get_secure_value "REDIS_PORT" "REDIS_PORT" 0 "Redis Port")
APP_ENV=$(get_secure_value "ENV" "NODE_ENV" 0 "Application Environment")
APP_PORT=$(get_secure_value "PORT" "APP_PORT" 0 "Application Port")

# PostgreSQL password para el secret separado
POSTGRES_PASSWORD=$DB_PASSWORD  # Usar la misma que PASS_DB

echo ""
echo "🔒 Encoding to base64..."

# Codificar a base64
JWT_SECRET_B64=$(echo -n "$JWT_SECRET" | base64 | tr -d '\n')
DB_PASSWORD_B64=$(echo -n "$DB_PASSWORD" | base64 | tr -d '\n')
REDIS_PASSWORD_B64=$(echo -n "$REDIS_PASSWORD" | base64 | tr -d '\n')
JWT_SALT_B64=$(echo -n "$JWT_SALT" | base64 | tr -d '\n')
POSTGRES_PASSWORD_B64=$(echo -n "$POSTGRES_PASSWORD" | base64 | tr -d '\n')

# Codificar hashes a base64
JWT_SECRET_HASH_B64=$(echo -n "$JWT_SECRET_HASH" | base64 | tr -d '\n')
DB_PASSWORD_HASH_B64=$(echo -n "$DB_PASSWORD_HASH" | base64 | tr -d '\n')
REDIS_PASSWORD_HASH_B64=$(echo -n "$REDIS_PASSWORD_HASH" | base64 | tr -d '\n')
JWT_SALT_HASH_B64=$(echo -n "$JWT_SALT_HASH" | base64 | tr -d '\n')

echo ""
echo "📁 Creating secret files..."

# Crear archivo de secrets principal
cat > manifests/shared/secrets.yaml << EOF
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
  labels:
    app: go-api
    generated-by: script
    created: $(date -Iseconds)
    source: env-file
type: Opaque
data:
  # Application secrets (base64 encoded)
  DB_PASSWORD: $DB_PASSWORD_B64
  REDIS_PASSWORD: $REDIS_PASSWORD_B64
  JWT_SECRET: $JWT_SECRET_B64
  JWT_SALT: $JWT_SALT_B64
  
  # Verification hashes (SHA256 in base64)
  DB_PASSWORD_HASH: $DB_PASSWORD_HASH_B64
  REDIS_PASSWORD_HASH: $REDIS_PASSWORD_HASH_B64
  JWT_SECRET_HASH: $JWT_SECRET_HASH_B64
  JWT_SALT_HASH: $JWT_SALT_HASH_B64
EOF

# Crear secret de PostgreSQL
cat > manifests/postgres/secret.yaml << EOF
apiVersion: v1
kind: Secret
metadata:
  name: postgres-secret
  labels:
    app: postgres
    generated-by: script
    created: $(date -Iseconds)
    source: env-file
type: Opaque
data:
  postgres-password: $POSTGRES_PASSWORD_B64
EOF

echo ""
echo "📄 Creating ConfigMap with your .env settings..."

# Crear ConfigMap con todas las configuraciones
cat > manifests/go-api/configmap.yaml << EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: go-api-config
  labels:
    app: go-api
    source: env-file
data:
  # Application Configuration
  APP_NAME: "SoftCalFut API"
  APP_ENV: "$APP_ENV"
  APP_PORT: "$APP_PORT"
  
  # Server Configuration
  SERVER_HOST: "0.0.0.0"
  SERVER_PORT: "$APP_PORT"
  
  # Database Configuration
  DB_HOST: "$DB_HOST"
  DB_PORT: "$DB_PORT"
  DB_USER: "$DB_USER"
  DB_NAME: "$DB_NAME"
  
  # Redis Configuration
  REDIS_HOST: "$REDIS_HOST"
  REDIS_PORT: "$REDIS_PORT"
  
  # JWT Configuration
  JWT_EXPIRES_IN: "24h"
  
  # Application Settings
  LOG_LEVEL: "info"
  CORS_ALLOWED_ORIGINS: "*"
EOF

echo ""
echo "📝 Creating deployment verification file..."

# Crear archivo de verificación
cat > manifests/shared/secrets-verification.txt << EOF
# 🔐 SECRETS VERIFICATION FILE
# Generated: $(date)
# Source: .env file
# 
# THIS FILE CONTAINS SENSITIVE INFORMATION
# DO NOT COMMIT TO VERSION CONTROL

## Application Configuration
APP_ENV: $APP_ENV
APP_PORT: $APP_PORT

## Database Configuration
DB_HOST: $DB_HOST
DB_PORT: $DB_PORT  
DB_USER: $DB_USER
DB_NAME: $DB_NAME
DB_PASSWORD: 
  Length: ${#DB_PASSWORD}
  Hash: $DB_PASSWORD_HASH

## Redis Configuration
REDIS_HOST: $REDIS_HOST
REDIS_PORT: $REDIS_PORT
REDIS_PASSWORD:
  Length: ${#REDIS_PASSWORD}
  Hash: $REDIS_PASSWORD_HASH

## JWT Configuration  
JWT_SECRET:
  Length: ${#JWT_SECRET}
  Hash: $JWT_SECRET_HASH
JWT_SALT:
  Length: ${#JWT_SALT} 
  Hash: $JWT_SALT_HASH

## Source Variables Used:
$(env | grep -E "(ENV|PORT|DB_|REDIS_|JWT_)" | sed 's/^/  /')

## Verification Command:
# To verify a secret in your app, compare the SHA256 hash:
# echo -n "your-secret-value" | sha256sum
EOF

echo "✅ Secrets generation completed!"
echo ""
echo "📋 Summary:"
echo "   ✅ manifests/shared/secrets.yaml (from your .env)"
echo "   ✅ manifests/postgres/secret.yaml (from your .env)" 
echo "   ✅ manifests/go-api/configmap.yaml (updated with .env settings)"
echo "   ✅ manifests/shared/secrets-verification.txt (DO NOT COMMIT)"
echo ""
echo "🔍 Your .env configuration has been converted to Kubernetes Secrets!"
echo "   DB Password: ${#DB_PASSWORD} chars"
echo "   Redis Password: ${#REDIS_PASSWORD} chars" 
echo "   JWT Secret: ${#JWT_SECRET} chars"
echo ""
echo "🚀 To deploy: make deploy"