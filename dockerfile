# syntax=docker/dockerfile:1.4

# Fase de construcción
FROM node:21-alpine3.19 AS builder

WORKDIR /usr/src/app

# 1. Copia solo los archivos necesarios para dependencias primero
COPY package.json package-lock.json .npmrc* ./

# 2. Instalación con caché de npm
RUN target=/root/.npm \
    npm ci --include=dev

# 3. Copia el resto de archivos de configuración
COPY tsconfig.json ./
COPY prisma ./prisma
RUN npx prisma generate

# 4. Copia el código fuente
COPY . .

# 5. Build con caché
RUN target=/root/.npm \
    npm run build

# Fase de producción final
FROM node:21-alpine3.19

WORKDIR /usr/src/app

# Instala PM2 y dependencias globales
RUN npm install -g pm2@latest

# Asegura que el home del usuario node y la carpeta .pm2 existen y son suyos
RUN mkdir -p /home/node/.pm2 && chown -R node:node /home/node

# Establece la variable de entorno HOME correctamente
ENV HOME=/home/node

# Copia desde builder (con ownership)
COPY --from=builder --chown=node:node /usr/src/app/package*.json ./ 
COPY --from=builder --chown=node:node /usr/src/app/node_modules ./node_modules
COPY --from=builder --chown=node:node /usr/src/app/dist ./dist
COPY --from=builder --chown=node:node /usr/src/app/prisma ./prisma
COPY --from=builder --chown=node:node /usr/src/app/ecosystem.config.js ./ecosystem.config.js
COPY --chown=node:node .env ./ 

# Cambia al usuario no root
RUN chown -R node:node /usr/src/app
    
USER node

EXPOSE 3000

# Ejecuta PM2 en modo producción sin daemon
CMD ["pm2-runtime", "ecosystem.config.js"]