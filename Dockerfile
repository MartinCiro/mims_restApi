# SOLO UNA ETAPA - Para desarrollo con hot-reload
FROM maven:3.9.11-eclipse-temurin-25

# Instalar herramientas útiles para desarrollo
RUN apt-get update && apt-get install -y \
    curl \
    wget \
    git \
    vim \
    && rm -rf /var/lib/apt/lists/*

# Configurar Maven para desarrollo
ENV MAVEN_CONFIG=/root/.m2

# Crear directorio de trabajo
WORKDIR /app

# Copiar SOLO el pom.xml primero (para cache de dependencias)
COPY pom.xml .

# Descargar dependencias (caché en capa separada)
RUN mvn dependency:go-offline

# NO copies el src aquí, se montará como volumen
# COPY src ./src  <-- ELIMINA ESTA LÍNEA

# Exponer puertos
EXPOSE 8080
EXPOSE 5006

# Comando por defecto - con devtools para hot reload
CMD ["mvn", "spring-boot:run", "-Dspring-boot.run.jvmArguments=-agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=*:5006"]