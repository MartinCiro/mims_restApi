# Etapa de construcción
FROM maven:3.9.9-eclipse-temurin-21 AS builder

# Instalar herramientas útiles para desarrollo
RUN apt-get update && apt-get install -y \
    curl \
    wget \
    git \
    vim \
    && rm -rf /var/lib/apt/lists/*

# Configurar Maven para desarrollo
ENV MAVEN_OPTS="-XX:+TieredCompilation -XX:TieredStopAtLevel=1"
ENV MAVEN_CONFIG=/root/.m2

# Crear directorio de trabajo
WORKDIR /app

# Copiar archivos del proyecto
COPY pom.xml .
COPY src ./src

# Descargar dependencias (caché en capa separada)
RUN mvn dependency:go-offline

# Compilar la aplicación
RUN mvn clean package -DskipTests

# Etapa final de ejecución
FROM eclipse-temurin:21-jre-alpine

# Instalar curl para health checks
RUN apk add --no-cache curl

# Crear usuario no-root para seguridad
RUN addgroup -S spring && adduser -S spring -G spring
USER spring:spring

WORKDIR /app

# Copiar el JAR desde la etapa de construcción
COPY --from=builder /app/target/*.jar app.jar

# Variables de entorno
ENV JAVA_OPTS=""
ENV SPRING_PROFILES_ACTIVE="docker"

# Exponer puerto
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8080/actuator/health || exit 1

# Punto de entrada
ENTRYPOINT ["sh", "-c", "java $JAVA_OPTS -jar /app/app.jar"]