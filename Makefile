# Makefile
.PHONY: build build-local deploy deploy-dry-run logs clean status \
        test configmaps secrets network-policies volumes services \
        port-forward exec-shell describe debug help generate-secrets verify-secrets deploy-secure

# Variables
APP_NAME ?= go-api
APP_IMAGE ?= $(APP_NAME):latest
K8S_NAMESPACE ?= default
MANIFESTS_DIR ?= manifests

## --- Desarrollo Local ---
build:
	@echo "🔨 Building Docker image..."
	-docker rmi $(APP_IMAGE) 2>/dev/null || true
	docker build -t $(APP_IMAGE) .

build-local: build
	@echo "🚀 Running locally..."
	docker run -p 4000:4000 --env-file .env $(APP_IMAGE)  # ← PUERTO 4000

test:
	@echo "🧪 Running tests..."
	go test ./... -v

## --- Secrets Management ---
generate-secrets:
	@echo "🔐 Generating secure secrets..."
	@chmod +x scripts/generate-secrets.sh
	@./scripts/generate-secrets.sh

verify-secrets:
	@echo "🔍 Verifying secrets..."
	@chmod +x scripts/verify-secrets.sh
	@./scripts/verify-secrets.sh

deploy-secure: generate-secrets deploy

## --- Kubernetes Deployment ---
deploy-dry-run:
	@echo "🔍 Dry-run deployment..."
	@echo "=== Volumes ==="
	kubectl apply -f $(MANIFESTS_DIR)/shared/volumes.yaml --dry-run=client -o yaml
	@echo "=== Secrets ==="
	kubectl apply -f $(MANIFESTS_DIR)/shared/secrets.yaml --dry-run=client -o yaml
	kubectl apply -f $(MANIFESTS_DIR)/postgres/secret.yaml --dry-run=client -o yaml
	@echo "=== PostgreSQL ==="
	kubectl apply -f $(MANIFESTS_DIR)/postgres/ --dry-run=client -o yaml
	@echo "=== Redis ==="
	kubectl apply -f $(MANIFESTS_DIR)/redis/ --dry-run=client -o yaml
	@echo "=== ConfigMaps ==="
	kubectl apply -f $(MANIFESTS_DIR)/go-api/configmap.yaml --dry-run=client -o yaml
	@echo "=== Go API ==="
	kubectl apply -f $(MANIFESTS_DIR)/go-api/ --dry-run=client -o yaml
	@echo "=== Network Policies ==="
	kubectl apply -f $(MANIFESTS_DIR)/shared/network-policies.yaml --dry-run=client -o yaml

deploy:
	@echo "🚀 Deploying to Kubernetes..."
	@echo "=== Creating volumes... ==="
	kubectl apply -f $(MANIFESTS_DIR)/shared/
	@echo "=== Creating secrets... ==="
	kubectl apply -f $(MANIFESTS_DIR)/shared/secrets.yaml
	kubectl apply -f $(MANIFESTS_DIR)/postgres/secret.yaml
	@echo "=== Deploying PostgreSQL... ==="
	kubectl apply -f $(MANIFESTS_DIR)/postgres/
	@echo "=== Deploying Redis... ==="
	kubectl apply -f $(MANIFESTS_DIR)/redis/
	@echo "=== Creating configmaps... ==="
	kubectl apply -f $(MANIFESTS_DIR)/go-api/configmap.yaml
	@echo "=== Deploying Go API... ==="
	kubectl apply -f $(MANIFESTS_DIR)/go-api/
	@echo "=== Applying network policies... ==="
	kubectl apply -f $(MANIFESTS_DIR)/shared/network-policies.yaml
	@echo "✅ Deployment completed!"
	@timeout 5 >nul 2>&1 || timeout 5 >nul 2>&1 || ping -n 6 127.0.0.1 >nul 2>&1 || sleep 5
	@make status

## --- Gestión de Recursos Específicos ---
configmaps:
	@echo "📄 ConfigMaps:"
	kubectl get configmaps

secrets:
	@echo "🔐 Secrets:"
	kubectl get secrets

network-policies:
	@echo "🛡️ Network Policies:"
	kubectl get networkpolicies

volumes:
	@echo "💾 Persistent Volumes:"
	kubectl get pv,pvc

services:
	@echo "🔌 Services:"
	kubectl get services

## --- Monitoreo y Debugging ---
logs:
	@echo "📋 Go API logs:"
	kubectl logs -l app=go-api -f

logs-postgres:
	@echo "📋 PostgreSQL logs:"
	kubectl logs -l app=postgres -f

logs-redis:
	@echo "📋 Redis logs:"
	kubectl logs -l app=redis -f

status:
	@echo "📊 Cluster Status:"
	@echo "=== Pods ==="
	kubectl get pods -o wide
	@echo "=== Services ==="
	kubectl get services
	@echo "=== Persistent Volumes ==="
	kubectl get pvc
	@echo "=== Network Policies ==="
	kubectl get networkpolicies

describe:
	@echo "🔍 Detailed status:"
	@echo "=== Go API ==="
	kubectl describe deployment go-api-deployment
	@echo "=== PostgreSQL ==="
	kubectl describe deployment postgres-deployment
	@echo "=== Redis ==="
	kubectl describe deployment redis-deployment

port-forward:
	@echo "🔗 Port forwarding Go API to localhost:4000..."  # ← PUERTO 4000
	kubectl port-forward svc/go-api-service 4000:4000       # ← PUERTO 4000

port-forward-db:
	@echo "🔗 Port forwarding PostgreSQL to localhost:5432..."
	kubectl port-forward svc/postgres-service 5432:5432

exec-shell:
	@echo "🐚 Accessing Go API pod shell..."
	@POD=$$(kubectl get pods -l app=go-api -o jsonpath='{.items[0].metadata.name}') && \
	kubectl exec -it $$POD -- /bin/sh

debug:
	@echo "🐛 Debug information:"
	@echo "=== Events ==="
	kubectl get events --sort-by=.metadata.creationTimestamp
	@echo "=== Resource usage ==="
	-kubectl top pods
	@echo "=== Go API environment ==="
	@POD=$$(kubectl get pods -l app=go-api -o jsonpath='{.items[0].metadata.name}') && \
	kubectl exec $$POD -- env | sort

## --- Limpieza ---
clean:
	@echo "🧹 Cleaning up resources..."
	@echo "=== Removing Go API... ==="
	-kubectl delete -f $(MANIFESTS_DIR)/go-api/
	@echo "=== Removing Redis... ==="
	-kubectl delete -f $(MANIFESTS_DIR)/redis/
	@echo "=== Removing PostgreSQL... ==="
	-kubectl delete -f $(MANIFESTS_DIR)/postgres/
	@echo "=== Removing Network Policies... ==="
	-kubectl delete -f $(MANIFESTS_DIR)/shared/network-policies.yaml
	@echo "=== Removing ConfigMaps... ==="
	-kubectl delete -f $(MANIFESTS_DIR)/go-api/configmap.yaml
	@echo "=== Removing Secrets... ==="
	-kubectl delete -f $(MANIFESTS_DIR)/shared/secrets.yaml
	-kubectl delete -f $(MANIFESTS_DIR)/postgres/secret.yaml
	@echo "=== Removing Volumes... ==="
	-kubectl delete -f $(MANIFESTS_DIR)/shared/volumes.yaml
	@echo "✅ Cleanup completed!"

clean-hard:
	@echo "💥 Hard cleanup (deleting namespace)..."
	-kubectl delete namespace $(K8S_NAMESPACE)
	@echo "✅ Hard cleanup completed!"

## --- Utilidades ---
restart:
	@echo "🔄 Restarting Go API deployment..."
	kubectl rollout restart deployment/go-api-deployment
	@make status

scale:
	@echo "📈 Scaling Go API to 3 replicas..."
	kubectl scale deployment go-api-deployment --replicas=3

scale-down:
	@echo "📉 Scaling Go API to 1 replica..."
	kubectl scale deployment go-api-deployment --replicas=1

health-check:
	@echo "❤️ Health checking..."
	@echo "=== Liveness ==="
	-kubectl get pods -l app=go-api -o jsonpath='{.items[*].status.conditions[?(@.type=="Ready")].status}'
	@echo "=== Readiness ==="
	@POD=$$(kubectl get pods -l app=go-api -o jsonpath='{.items[0].metadata.name}') && \
	kubectl exec $$POD -- wget -q -O- http://localhost:4000/health  # ← PUERTO 4000

## --- NUEVO COMANDO: ACCESO NODEPORT ---
nodeport-access:
	@echo "🌐 NodePort Access Information:"
	@echo "=== Service Info ==="
	@kubectl get service go-api-service -o jsonpath='{range .spec.ports[0]}{.nodePort}{end}' | xargs -I {} echo "NodePort: {}"
	@echo "=== Node IPs ==="
	@kubectl get nodes -o jsonpath='{.items[*].status.addresses[?(@.type=="ExternalIP")].address}'
	@echo ""
	@echo "📝 Access URLs:"
	@NODE_PORT=$$(kubectl get service go-api-service -o jsonpath='{range .spec.ports[0]}{.nodePort}{end}') && \
	NODE_IP=$$(kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type=="ExternalIP")].address}') && \
	if [ -z "$$NODE_IP" ]; then \
		NODE_IP=$$(kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type=="InternalIP")].address}'); \
	fi && \
	echo "   NodePort: http://$$NODE_IP:$$NODE_PORT" && \
	echo "   Port-Forward: make port-forward (then http://localhost:4000)"

## --- Ayuda ---
help:
	@echo "🚀 Available commands:"
	@echo ""
	@echo "🏗️  Development:"
	@echo "  build               Build Docker image"
	@echo "  build-local         Build and run locally (port 4000)"
	@echo "  test                Run tests"
	@echo ""
	@echo "🔐 Secrets Management:"
	@echo "  generate-secrets    Generate Kubernetes secrets from .env"
	@echo "  verify-secrets      Verify generated secrets"
	@echo "  deploy-secure       Generate secrets and deploy"
	@echo ""
	@echo "🚀 Deployment:"
	@echo "  deploy-dry-run      Preview deployment without applying"
	@echo "  deploy              Full deployment to Kubernetes"
	@echo ""
	@echo "📊 Monitoring:"
	@echo "  status              Show cluster status"
	@echo "  logs                Tail Go API logs"
	@echo "  logs-postgres       Tail PostgreSQL logs"
	@echo "  logs-redis          Tail Redis logs"
	@echo "  describe            Detailed resource information"
	@echo "  debug               Debug information"
	@echo ""
	@echo "🔧 Utilities:"
	@echo "  port-forward        Forward Go API to localhost:4000"  # ← ACTUALIZADO
	@echo "  port-forward-db     Forward PostgreSQL to localhost:5432"
	@echo "  exec-shell          Access Go API pod shell"
	@echo "  nodeport-access     Show NodePort access information"  # ← NUEVO
	@echo "  restart             Restart Go API deployment"
	@echo "  scale               Scale Go API to 3 replicas"
	@echo "  scale-down          Scale Go API to 1 replica"
	@echo "  health-check        Health check endpoints (port 4000)" # ← ACTUALIZADO
	@echo ""
	@echo "🧹 Cleanup:"
	@echo "  clean               Remove all resources"
	@echo "  clean-hard          Hard cleanup (delete namespace)"
	@echo ""
	@echo "📁 Specific Resources:"
	@echo "  configmaps          List ConfigMaps"
	@echo "  secrets             List Secrets"
	@echo "  network-policies    List Network Policies"
	@echo "  volumes             List volumes"
	@echo "  services            List services"
	@echo ""
	@echo "❓ Help:"
	@echo "  help                Show this help message"

# Default target
.DEFAULT_GOAL := help