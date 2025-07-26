---

# 📌 **Guía para Configurar un Contenedor Postgres con Docker**

Este documento explica cómo configurar un contenedor de **Postgres** con Docker, realizar operaciones con la base de datos y entender la estructura del proyecto.
---


## 🚀 **1. Iniciar el Contenedor**  

Para desplegar el contenedor, ejecuta (renombrar example a .env sino lo tiene []):

```bash
docker-compose up -d
```

---

## 🧪 **2. Ejecutar Pruebas**

Para iniciar las pruebas de microservicios, usa:

```bash
docker exec -it nestjs_app sh -c "npm test"
```

---

## 📄 **3. Acceder a la Documentación de la API**

La documentación de la API se aloja en la siguiente ruta (el puerto varia segun el .env):

```
http://localhost:4000/api/docs
```

Si la documentación no carga, verifica que el archivo `docs.html` exista en la carpeta `public/`.

---

# 🏗 **Configuración de la Base de Datos**

## **4. Crear el Contenedor de PostgreSQL**

Para desplegar un contenedor con **PostgreSQL** y exponer el puerto `5432`, ejecuta:

```bash
docker run --name psql -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=your_password -e POSTGRES_DB=nom_bd_psql -v postgres_data:/var/lib/postgresql/data -p 5432:5432 -d postgres:16-alpine
```

---

## **5. Crear la Base de Datos**  

Si la base de datos `nom_bd_psql` no existe, créala con:

```bash
docker exec -it psql psql -U postgres -c "CREATE DATABASE nom_bd_psql;"
```

---

## **6. Crear una Copia de Seguridad**  

Para generar un respaldo de la base de datos `nom_bd_psql`, usa:

```bash
docker exec -t psql bash -c 'pg_dump -U postgres -Fc -d nom_bd_psql > /tmp/backup.dump'
docker cp psql:/tmp/backup.dump ./src/bd_backup/backup.dump

```
---

## **7. Restaurar una Copia de Seguridad**

Para restaurar una copia de seguridad, ejecuta:

```bash
docker cp ./src/bd_backup/backup.dump psql:/tmp/backup.dump
docker exec -it psql pg_restore -U postgres -d nom_bd_psql -Fc /tmp/backup.dump

```

---
## **8. Variables de Entorno****

Verifique o renombre el archivo example, el siguiente comando se usa para generar el JWT_SECRETO

```bash
openssl rand -hex 32

```

---

# 📂 **Estructura del Back**

Este proyecto sigue una organización modular basada en **arquitectura hexagonal**, lo que mejora la mantenibilidad y escalabilidad.

```
/src
  ├── /bd_backup             # Carpeta para copias de seguridad
  ├── /core                  # Lógica de negocio
  │   ├── /usuarios          # Módulo de usuarios
  │   │   ├── usuarioPort.ts     # Interfaz del servicio
  │   │   ├── usuario.service.ts # Lógica de negocio de usuarios
  │
  ├── /interfaces            # Interfaces de comunicación
  │   ├── /api               # API REST
  │   │   ├── /usuarios      # Endpoints de usuarios
  │   │   │   ├── usuario.controller.ts  # Controlador de usuarios
  │   │   │   ├── usuario.module.ts      # Módulo de usuarios
  │   │   │   ├── /dtos                  # DTOs (Data Transfer Objects)
  │   │   │   │   ├── crear-usuario.dto.ts  # DTO para validación
  │   │   ├── /models            # Modelos de datos compartidos
  │   │   │   ├── ResponseBody.ts # Estructura estándar de respuesta
  │
  ├── /config               # Configuración global (NestJS, variables de entorno)
  ├── /shared               # Código reutilizable (utilidades, middlewares)
  ├── app.module.ts         # Módulo raíz de la aplicación
  ├── main.ts               # Punto de entrada de la aplicación (bootstrap de NestJS)
```

---

# 📌 **Descripción General**

- **`/core`** → Contiene la lógica de negocio pura, independiente de la infraestructura.  
- **`/interfaces/api`** → Define los controladores, módulos y DTOs para la API REST.  
- **`/config`** → Configuración general del proyecto (variables de entorno, NestJS, etc.).  
- **`/shared`** → Código reutilizable (utilidades, excepciones, middlewares).  

Esta estructura modular **facilita la escalabilidad y el mantenimiento** del código. 🚀  

---