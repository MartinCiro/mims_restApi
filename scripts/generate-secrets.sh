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
    
    # Eliminar todos los caracteres whitespace (espacios, tabs, nuevas líneas)
    value=$(printf '%s' "$value" | tr -d '[:space:]')
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
        echo "   ✅ Using $primary_var: $description" >&2
        value="${!primary_var}"
    elif [ -n "${!fallback_var}" ]; then
        echo "   ✅ Using $fallback_var (fallback): $description" >&2
        value="${!fallback_var}"
    else
        echo "   ⚠️  Neither $primary_var nor $fallback_var set, generating random value" >&2
        value=$(openssl rand -base64 "$default_length" | head -c "$default_length")
    fi
    
    # Limpiar espacios, tabs y nuevas líneas
    value=$(echo "$value" | tr -d '[:space:]')
    
    # Generar SHA256 del valor (ya limpio)
    local hash
    hash=$(echo -n "$value" | sha256sum | cut -d' ' -f1)
    
    # Devolver SOLO valor y hash (sin mensajes)
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

# Obtener valores para SWAG
DUCKDNS_TOKEN=$(get_secure_value "DUCKDNS_TOKEN" "DUCKDNS_TOKEN" 0 "DuckDNS Token")
SWAG_EMAIL=$(get_secure_value "EMAIL" "SWAG_EMAIL" 0 "SWAG Email")
SWAG_DOMAIN=$(get_secure_value "DOMINIO" "SWAG_DOMAIN" 0 "SWAG Domain")

# Obtener otras configuraciones importantes
USER_DB=$(get_secure_value "USER_DB" "DB_USER" 0 "Database User")
NAME_DB=$(get_secure_value "NAME_DB" "DB_NAME" 0 "Database Name")
HOST_DB=$(get_secure_value "HOST_DB" "DB_HOST" 0 "Database Host")
PORT_DB=$(get_secure_value "PORT_DB" "DB_PORT" 0 "Database Port")
REDIS_URL=$(get_secure_value "REDIS_URL" "REDIS_HOST" 0 "Redis Host")
REDIS_PORT=$(get_secure_value "REDIS_PORT" "REDIS_PORT" 0 "Redis Port")
APP_ENV=$(get_secure_value "Env" "NODE_ENV" 0 "Application Environment")
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

# Codificar secrets de SWAG a base64
DUCKDNS_TOKEN_B64=$(echo -n "$DUCKDNS_TOKEN" | base64 | tr -d '\n')
SWAG_EMAIL_B64=$(echo -n "$SWAG_EMAIL" | base64 | tr -d '\n')

echo ""
echo "📁 Creating secret files..."

# Crear archivo de secrets principal
cat > manifests/shared/secrets.yaml << EOF
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets-1
  labels:
    app: go-api
    generated-by: script
    created: $(date -I)
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
  name: postgres-secret-1
  labels:
    app: postgres
    generated-by: script
    created: $(date -I)
    source: env-file
type: Opaque
data:
  postgres-password: $POSTGRES_PASSWORD_B64
EOF

# Crear secret de SWAG
cat > manifests/swag/secret.yaml << EOF
apiVersion: v1
kind: Secret
metadata:
  name: swag-secrets-1
  labels:
    app: swag
    generated-by: script
    created: $(date -I)
    source: env-file
type: Opaque
data:
  duckdns-token: $DUCKDNS_TOKEN_B64
  email: $SWAG_EMAIL_B64
EOF

echo ""
echo "📄 Creating ConfigMap with your .env settings..."

# Crear ConfigMap con todas las configuraciones
cat > manifests/go-api/configmap.yaml << EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: go-api-config-1
  labels:
    app: go-api
    source: env-file
    instance: "1"
data:
  # Application Configuration
  APP_NAME: "Template API - Instance 1"
  APP_ENV: "$APP_ENV"
  APP_PORT: "$APP_PORT"
  
  # Server Configuration
  SERVER_HOST: "0.0.0.0"
  SERVER_PORT: "$APP_PORT"
  
  # Database Configuration
  HOST_DB: "$HOST_DB"
  PORT_DB: "$PORT_DB"
  USER_DB: "$USER_DB"
  NAME_DB: "$NAME_DB"
  
  # Redis Configuration
  REDIS_URL: "$REDIS_URL"
  REDIS_PORT: "$REDIS_PORT"
  
  # JWT Configuration
  JWT_EXPIRES_IN: "24h"
  
  # Application Settings
  LOG_LEVEL: "info"
  CORS_ALLOWED_ORIGINS: "*"
EOF

# Crear ConfigMap para SWAG
cat > manifests/swag/configmap.yaml << EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: swag-proxy-configs-1
  labels:
    app: swag
    generated-by: script
    instance: "1"
data:
  http-challenge.conf: |
    server {
        listen 80;
        listen [::]:80;
        server_name $SWAG_DOMAIN.duckdns.org *.$SWAG_DOMAIN.duckdns.org;
        
        # Challenge de Let's Encrypt
        location /.well-known/acme-challenge/ {
            root /config/www;
            try_files \$uri =404;
        }
        
        # Todo lo demás redirige a HTTPS
        location / {
            return 301 https://\$host\$request_uri;
        }
    }
    
  # Proxy para tu API Go
  go-api.subdomain.conf: |
    server {
        listen 443 ssl;
        listen [::]:443 ssl;
        server_name api.$SWAG_DOMAIN.duckdns.org;
        
        include /config/nginx/ssl.conf;
        
        client_max_body_size 0;

        # Forzar HTTP/1.1 para evitar problemas con HTTP/2
        proxy_http_version 1.1;
        
        # Configuración específica para la raíz
        location / {
            proxy_pass http://go-api-service-1:4000/;
            proxy_set_header Host \$host;
            proxy_set_header X-Real-IP \$remote_addr;
            proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto \$scheme;
        }

        # Health check para la API
        location /health {
            proxy_pass http://go-api-service-1:4000/health;
            proxy_set_header Host \$host;
        }
    }

  www.subdomain.conf: |
    server {
        listen 443 ssl;
        listen [::]:443 ssl;
        server_name www.$SWAG_DOMAIN.duckdns.org;
        
        include /config/nginx/ssl.conf;

        return 301 https://$SWAG_DOMAIN.duckdns.org\$request_uri;
    }
  
  
  default.conf: |
    server {
        listen 443 ssl;
        listen [::]:443 ssl;
        server_name $SWAG_DOMAIN.duckdns.org;
        
        include /config/nginx/ssl.conf;
        
        # Redirige el dominio principal a www
        return 301 https://www.$SWAG_DOMAIN.duckdns.org\$request_uri;
    }
    server {
        listen 443 ssl default_server;
        listen [::]:443 ssl default_server;
        server_name _;
        
        include /config/nginx/ssl.conf;
        
        # Redirige cualquier subdominio no configurado a www
        return 301 https://www.$SWAG_DOMAIN.duckdns.org\$request_uri;
    }
    
    include /config/nginx/proxy-confs/*.subdomain.conf;
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
HOST_DB: $HOST_DB
PORT_DB: $PORT_DB  
USER_DB: $USER_DB
NAME_DB: $NAME_DB
DB_PASSWORD: 
  Length: ${#DB_PASSWORD}
  Hash: $DB_PASSWORD_HASH

## Redis Configuration
REDIS_URL: $REDIS_URL
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

## SWAG Configuration
SWAG_DOMAIN: $SWAG_DOMAIN
SWAG_EMAIL: $SWAG_EMAIL
DUCKDNS_TOKEN:
  Length: ${#DUCKDNS_TOKEN}

## Source Variables Used:
$(env | grep -E "(ENV|PORT|DB_|REDIS_|JWT_|DOMINIO|DUCKDNS|EMAIL)" | sed 's/^/  /')

## Verification Command:
# To verify a secret in your app, compare the SHA256 hash:
# echo -n "your-secret-value" | sha256sum
EOF

echo "✅ Secrets generation completed!"
echo ""
echo "📋 Summary:"
echo "   ✅ manifests/shared/secrets.yaml"
echo "   ✅ manifests/postgres/secret.yaml" 
echo "   ✅ manifests/swag/secret.yaml"
echo "   ✅ manifests/swag/configmap.yaml"
echo "   ✅ manifests/go-api/configmap.yaml"
echo "   ✅ manifests/shared/secrets-verification.txt"
echo ""
echo "🔍 Your .env configuration has been converted to Kubernetes Secrets!"
echo "   DB Password: ${#DB_PASSWORD} chars"
echo "   Redis Password: ${#REDIS_PASSWORD} chars" 
echo "   JWT Secret: ${#JWT_SECRET} chars"
echo "   SWAG Domain: $SWAG_DOMAIN.duckdns.org"
echo ""
echo "🚀 To deploy: make deploy"