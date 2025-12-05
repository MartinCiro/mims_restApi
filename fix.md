# 📋 **README: Resolución de Problemas de Conectividad API Go en Kubernetes**

## 🎯 **Problemas Resueltos**

### **1. Problema de Firewall** 🔥
**Síntoma:** API inaccesible desde internet, error "connection refused"
**Causa:** Firewall bloqueando puertos 80 y 443
**Solución:** 
```bash
sudo firewall-cmd --permanent --add-service=http --add-service=https
sudo firewall-cmd --permanent --add-port=8080/tcp --add-port=6443/tcp
sudo firewall-cmd --reload
```

### **2. Problema de DNS Inconsistente** 🌐
**Síntoma:** DuckDNS apuntando a múltiples IPs diferentes
**Causa:** DNS desincronizado entre servidores (Cloudflare vs Google vs Quad9)
**Solución:**
```bash
curl "https://www.duckdns.org/update?domains=iyata,api.iyata&token=TU_TOKEN&ip=149.56.131.106&ipv6=2607:5300:205:200::83ff"
```

### **3. Problema de Traefik IPv6-only** 🔌
**Síntoma:** Traefik solo escuchando en IPv6
**Causa:** Configuración de entrypoints incorrecta
**Solución:**
```bash
kubectl edit deployment -n traefik traefik
# Cambiar: --entrypoints.web.address=:80
# A:      --entrypoints.web.address=0.0.0.0:80
```

### **4. Problema de Providers Traefik** ⚙️
**Síntoma:** Traefik no procesaba Ingress regulares, solo CRDs
**Causa:** Falta de configuración `--providers.kubernetesingress`
**Solución:**
```bash
kubectl patch deployment -n traefik traefik -p '{"spec":{"template":{"spec":{"containers":[{"name":"traefik","args":[
  "--api.dashboard=true",
  "--api.insecure=true", 
  "--providers.kubernetescrd",
  "--providers.kubernetesingress",
  "--entrypoints.web.address=0.0.0.0:80",
  "--entrypoints.web.http2.maxconcurrentstreams=100",
  "--entrypoints.websecure.address=0.0.0.0:443",
  "--certificatesresolvers.letsencrypt.acme.email=martinciro11@gmail.com",
  "--certificatesresolvers.letsencrypt.acme.storage=/data/acme.json", 
  "--certificatesresolvers.letsencrypt.acme.httpchallenge.entrypoint=web"
]}]}}}}'
```

### **5. Problema de Reglas CRD Corruptas** 🐛
**Síntoma:** Caracteres de escape incorrectos en reglas Traefik
**Causa:** IngressRoute con reglas mal formateadas
**Solución:**
```bash
kubectl delete ingressroutes.traefik.containo.us go-api-ingress
```

### **6. Problema de Conflictos Cloudflare** ☁️
**Síntoma:** DNS apuntando a tunnel de Cloudflare en lugar de IP real
**Causa:** Tunnel de Cloudflare configurado
**Solución:** Eliminar o pausar tunnel en Cloudflare Dashboard

## 🛠️ **Comandos Clave y su Función**

### **Diagnóstico:**
```bash
# Verificar estado del firewall
sudo firewall-cmd --list-all

# Verificar puertos escuchando
sudo netstat -tulpn | grep -E ":80 |:443 "

# Verificar DNS desde múltiples servidores
nslookup api.iyata.duckdns.org 1.1.1.1
nslookup api.iyata.duckdns.org 8.8.8.8

# Verificar routers de Traefik
curl -s http://localhost:8080/api/http/routers | python3 -m json.tool
```

### **Configuración Firewall:**
```bash
# Activar y configurar firewalld
sudo systemctl enable firewalld
sudo systemctl start firewalld

# Abrir puertos esenciales
sudo firewall-cmd --permanent --add-service=http --add-service=https
sudo firewall-cmd --permanent --add-port=8080/tcp --add-port=6443/tcp
sudo firewall-cmd --reload
```

### **Kubernetes/Traefik:**
```bash
# Verificar Ingress
kubectl get ingress -o wide

# Verificar servicios
kubectl get services -A

# Reiniciar Traefik
kubectl rollout restart deployment -n traefik traefik

# Verificar logs
kubectl logs -n traefik -l app=traefik --tail=20
```

### **DuckDNS:**
```bash
# Actualizar IP IPv4 e IPv6
curl "https://www.duckdns.org/update?domains=iyata,api.iyata&token=TU_TOKEN&ip=149.56.131.106&ipv6=2607:5300:205:200::83ff"
```

## ✅ **Verificación Final**

### **Interno:**
```bash
# Desde el servidor
curl -v -H "Host: api.iyata.duckdns.org" http://149.56.131.106/health

# Desde dentro del cluster
kubectl exec -it go-api-deployment-1-6f7c6cd798-xw2mz -- wget -q -O - http://localhost:4000/health
```

### **Externo:**
```bash
# Desde cualquier computadora
curl -v http://api.iyata.duckdns.org/health

# Con HTTPS (usando IPv6)
curl -vk https://api.iyata.duckdns.org/health
```

## 📈 **Lecciones Aprendidas**

1. **Firewall primero:** Siempre verificar puertos abiertos antes de diagnosticar problemas complejos
2. **DNS inconsistente:** Los servidores DNS pueden tener cachés diferentes, verificar múltiples fuentes
3. **IPv4 vs IPv6:** Traefik debe configurarse explícitamente para escuchar en IPv4 (`0.0.0.0`)
4. **Providers Traefik:** Necesita ambos `kubernetescrd` y `kubernetesingress` para soportar todos los tipos de Ingress
5. **Propagación DNS:** Paciencia - puede tomar de 5 minutos a 2 horas para propagación completa
6. **Cloudflare conflicts:** Los tunnels de Cloudflare pueden interferir con DNS directo

## 🚀 **Estado Final**

✅ **API Go funcionando:** `http://api.iyata.duckdns.org/health`  
✅ **HTTP/HTTPS:** Accesible vía HTTP y HTTPS  
✅ **IPv4/IPv6:** Soporte dual-stack  
✅ **Firewall configurado:** Solo puertos necesarios abiertos  
✅ **DNS consistente:** DuckDNS apuntando a IP correcta  
✅ **Traefik funcionando:** Procesando Ingress y CRDs correctamente  

**Respuesta esperada:**
```json
{"ok":true,"statusCode":200,"result":{"message":"Hello world","status":"running","version":"1.0.0"}}
```

## 🆘 **Solución Rápida para Problemas Futuros**

1. **¿No responde?** → Verificar firewall: `sudo firewall-cmd --list-all`
2. **¿DNS incorrecto?** → Actualizar DuckDNS y esperar 5 minutos
3. **¿404 Not Found?** → Verificar Ingress: `kubectl get ingress`
4. **¿Error TLS?** → Verificar logs Traefik: `kubectl logs -n traefik -l app=traefik | grep -i acme`

**¡API lista para producción!** 🎉