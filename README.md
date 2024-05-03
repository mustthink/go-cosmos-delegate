# go-cosmos-delegate
simple application which parse transactions with delegate messages from cosmos blockchain 

## Installation
Clone the repository via `git clone`
> If you have Go installed
> install dependencies via `go mod download`

## Running

### Configuration
You can configure the application via environment variables.
- `CONFIG_PATH` - path to the configuration file.

You can configure the application via flags on the start:
- `--config` - path to the configuration file.

### 1) Run in Docker
1) Run PostgreSQL locally or via Docker.
>You can use `task docker-postgres` or `docker-compose up postgres` for running.
2) Run application via command `task docker-migrations` which runs migrations 
3) Run application via command `task docker-parser` which starts the application.

### 2) Run via Go
1) Run PostgreSQL locally or via Docker. (check details in the previous section)
2) Run migrations via `task migrate`
3) Run application via `task run`