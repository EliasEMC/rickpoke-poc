# Rick & Poke API

Este proyecto es una implementación de una API que combina datos de Rick and Morty y Pokémon, utilizando arquitectura hexagonal (también conocida como Ports and Adapters) para mantener una clara separación de responsabilidades y facilitar el mantenimiento y las pruebas.

## 🏗️ Arquitectura

El proyecto sigue los principios de la arquitectura hexagonal, dividiéndose en las siguientes capas:

### 1. Dominio (Core)
- `internal/domain/model`: Define las entidades del dominio (Rick, Pokemon)
- `internal/domain/service`: Define los puertos (interfaces) que el dominio necesita

### 2. Puertos (Ports)
- `internal/domain/service/character_repository.go`: Define la interfaz para el repositorio de personajes
- `internal/domain/service/rick_service.go`: Define la interfaz para el servicio de Rick and Morty
- `internal/domain/service/pokemon_service.go`: Define la interfaz para el servicio de Pokémon

### 3. Adaptadores (Adapters)
- `internal/adapter/api`: Implementaciones de los servicios externos (Rick and Morty API, Pokémon API)
- `internal/adapter/db`: Implementaciones del repositorio (PostgreSQL, MongoDB)
- `internal/adapter/handler`: Manejadores HTTP para las rutas

### 4. Casos de Uso (Application)
- `internal/usecase`: Implementa la lógica de negocio usando los puertos

### 5. Infraestructura
- `internal/infrastructure/http`: Configuración del router y middleware
- `internal/infrastructure/circuitbreaker`: Implementación del circuit breaker para las APIs externas

## 🚀 Configuración

### Variables de Entorno
Crea un archivo `.env` en la raíz del proyecto con las siguientes variables:

```env
# Server
PORT=8080                    # Puerto del servidor
LOG_LEVEL=debug             # Nivel de logging (debug, info, warn, error)
TIMEOUT=5s                  # Timeout para las peticiones HTTP

# External APIs
RICK_URL=https://rickandmortyapi.com/api    # URL base de la API de Rick and Morty
POKE_URL=https://pokeapi.co/api/v2          # URL base de la API de Pokémon

# Database Configuration
# Opciones: postgres o mongodb
DB_TYPE=postgres            # Tipo de base de datos a utilizar
DB_HOST=localhost          # Host de la base de datos
# Puerto: 5432 para postgres, 27017 para mongodb
DB_PORT=5432               # Puerto de la base de datos
# Usuario: alaska-eng para postgres, root para mongodb
DB_USER=                   # Usuario de la base de datos
DB_PASSWORD=               # Contraseña de la base de datos
DB_NAME=rickpoke          # Nombre de la base de datos
DB_SSL_MODE=disable       # Modo SSL para PostgreSQL
```

### Requisitos Previos

#### Para PostgreSQL:
```bash
# Instalar PostgreSQL
brew install postgresql@14

# Iniciar PostgreSQL
brew services start postgresql@14

# Crear la base de datos
createdb rickpoke
```

#### Para MongoDB:
```bash
# Instalar MongoDB
brew tap mongodb/brew
brew install mongodb-community

# Iniciar MongoDB
brew services start mongodb-community
```

## 🛠️ Instalación y Ejecución

1. Clonar el repositorio:
```bash
git clone <repository-url>
cd rickpoke-poc
```

2. Instalar dependencias:
```bash
go mod download
```

3. Ejecutar la aplicación:
```bash
make run
```

## 📡 Endpoints

### 1. Health Check
```bash
curl http://localhost:8080/health
```
Respuesta:
```json
{
  "status": "ok",
  "services": {
    "database": "ok",
    "rick_api": "ok",
    "poke_api": "ok"
  }
}
```

### 2. Guardar Personaje
```bash
curl -X POST http://localhost:8080/characters \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "name": "Rick Sanchez"
  }'
```

### 3. Listar Personajes
```bash
curl http://localhost:8080/characters
```

### 4. Obtener Personaje Específico
```bash
curl http://localhost:8080/character/1
```

### 5. Obtener Pokémon
```bash
curl http://localhost:8080/pokemon/pikachu
```

### 6. Obtener Datos Combinados
```bash
curl "http://localhost:8080/combined?char_id=1&pokemon=pikachu"
```

## 🔄 Flujo de Peticiones

1. **Guardar Personaje**:
   - La petición llega al `CharacterHandler`
   - El handler utiliza el caso de uso `StoreCharacter`
   - El caso de uso utiliza el repositorio configurado (PostgreSQL o MongoDB)
   - Los datos se almacenan en la base de datos

2. **Obtener Personaje**:
   - La petición llega al `CharacterHandler`
   - El handler utiliza el caso de uso `FetchCharacter`
   - El caso de uso utiliza el servicio de Rick and Morty a través del circuit breaker
   - Los datos se devuelven al cliente

3. **Obtener Pokémon**:
   - La petición llega al `PokemonHandler`
   - El handler utiliza el caso de uso `FetchPokemon`
   - El caso de uso utiliza el servicio de Pokémon a través del circuit breaker
   - Los datos se devuelven al cliente

4. **Obtener Datos Combinados**:
   - La petición llega al router
   - El router utiliza el caso de uso `CombinedUC`
   - El caso de uso combina datos de ambos servicios
   - Los datos combinados se devuelven al cliente

## 🛡️ Circuit Breaker

El proyecto implementa circuit breakers para las APIs externas:
- Se abre después de 5 fallos consecutivos
- Permite 3 peticiones en estado semi-abierto
- Timeout de 5 segundos antes de intentar recuperarse

Para probar el circuit breaker:
```bash
# Forzar fallos para abrir el circuito
for i in {1..10}; do
  curl http://localhost:8080/pokemon/nonexistentpokemon
  echo
  sleep 1
done

# Verificar estado
curl http://localhost:8080/health
```

## 🧪 Pruebas

Ejecutar las pruebas:
```bash
make test
```

Generar reporte de cobertura:
```bash
make coverage
```

## 📚 Estructura de Directorios

```
.
├── cmd/
│   └── server/
│       └── main.go           # Punto de entrada de la aplicación
├── internal/
│   ├── adapter/
│   │   ├── api/             # Adaptadores para APIs externas
│   │   ├── db/              # Adaptadores para bases de datos
│   │   └── handler/         # Manejadores HTTP
│   ├── domain/
│   │   ├── model/          # Entidades del dominio
│   │   └── service/        # Puertos (interfaces)
│   ├── infrastructure/
│   │   ├── circuitbreaker/ # Implementación del circuit breaker
│   │   └── http/           # Configuración HTTP
│   └── usecase/            # Casos de uso
├── pkg/
│   └── utils/              # Utilidades compartidas
├── .env                    # Variables de entorno
├── go.mod                  # Dependencias de Go
└── Makefile               # Comandos de construcción
```

## 🔄 Cambiar entre Bases de Datos

Para cambiar entre PostgreSQL y MongoDB:

1. Modifica el archivo `.env`:
   - Para PostgreSQL:
     ```env
     DB_TYPE=postgres
     DB_PORT=5432
     DB_USER=alaska-eng
     ```
   - Para MongoDB:
     ```env
     DB_TYPE=mongodb
     DB_PORT=27017
     DB_USER=root
     ```

2. Reinicia la aplicación:
   ```bash
   make run
   ```

## 📝 Notas Adicionales

- El proyecto utiliza Gin como framework web
- Implementa logging estructurado con Zap
- Usa circuit breakers para manejar fallos en APIs externas
- Soporta múltiples bases de datos mediante adaptadores
- Implementa health checks para monitoreo

## 🐳 Docker

El proyecto está configurado para ejecutarse con Docker y Docker Compose. Esto permite ejecutar la aplicación y sus dependencias (PostgreSQL y MongoDB) en contenedores aislados.

### Requisitos
- Docker
- Docker Compose

### Ejecutar con Docker

1. Construir y levantar los contenedores:
```bash
docker-compose up --build
```

2. Para ejecutar en segundo plano:
```bash
docker-compose up -d
```

3. Para detener los contenedores:
```bash
docker-compose down
```

4. Para ver los logs:
```bash
docker-compose logs -f
```

### Cambiar entre Bases de Datos

Para cambiar entre PostgreSQL y MongoDB, modifica el archivo `.env` en el contenedor de la aplicación:

1. Para PostgreSQL:
```env
DB_TYPE=postgres
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
```

2. Para MongoDB:
```env
DB_TYPE=mongodb
DB_HOST=mongodb
DB_PORT=27017
DB_USER=root
DB_PASSWORD=example
```

### Volúmenes

Los datos de las bases de datos se persisten en volúmenes Docker:
- `postgres_data`: Datos de PostgreSQL
- `mongodb_data`: Datos de MongoDB

### Migraciones

Las migraciones de PostgreSQL se ejecutan automáticamente al iniciar el contenedor por primera vez. Los scripts de migración se encuentran en el directorio `migrations/`. 