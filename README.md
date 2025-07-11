# dockenv

Simple local development environments with Docker Compose

## One-liner installation

```bash
curl -s https://raw.githubusercontent.com/mohammed-bageri/dockenv/main/install.sh | bash
```

## Or install manually (requires Go)

```bash
go install github.com/mohammed-bageri/dockenv@latest
```

dockenv is a CLI tool that helps you set up local development environments using Docker Compose. Install services like MySQL, PostgreSQL, Redis, MongoDB, and Kafka without installing them directly on your system.

[![Go Version](https://img.shields.io/badge/go-%3E%3D%201.24-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Release](https://img.shields.io/github/release/mohammed-bageri/dockenv.svg)](https://github.com/mohammed-bageri/dockenv/releases)

## Features

### Core Features (v0.1)

- 🚀 **Interactive Setup**: `dockenv init` with guided service selection
- 🐳 **Docker Integration**: Automatic Docker Compose file generation
- 📊 **Service Management**: Start, stop, restart services with simple commands
- 💾 **Persistent Data**: Volume management for data persistence
- 🔧 **Docker Validation**: Automatic Docker/Docker Compose installation checks
- 🎯 **One-liner Install**: `curl -s https://raw.githubusercontent.com/mohammed-bageri/dockenv/main/install.sh | bash`

### Advanced Features (v0.2)

- ➕ **Dynamic Services**: Add/remove services without restarting everything
- 🔄 **Auto-start**: System boot integration with systemd
- 🌍 **Environment Files**: Automatic `.env` file generation
- 🔍 **Project Detection**: Auto-detect Laravel, Node.js, Django projects
- 📋 **Profiles**: Pre-configured service bundles for common stacks
- 📝 **Status Monitoring**: Real-time service status and logs

## Quick Start

### Installation

```bash
# One-liner installation
curl -s https://raw.githubusercontent.com/mohammed-bageri/dockenv/main/install.sh | bash

# Or install manually (requires Go)
go install github.com/mohammed-bageri/dockenv@latest
```

### Basic Usage

```bash
# Initialize your development environment
dockenv init

# Start services
dockenv up

# Check status
dockenv status

# Stop services
dockenv down
```

### Profile-based Setup

```bash
# Laravel stack (MySQL + Redis)
dockenv init --profile laravel

# Node.js stack (PostgreSQL + Redis)
dockenv init --profile node

# Full stack (all services)
dockenv init --profile full
```

## Supported Services

| Service           | Description                   | Default Port | Environment Variables                                                            |
| ----------------- | ----------------------------- | ------------ | -------------------------------------------------------------------------------- |
| **MySQL**         | MySQL 8.0 Database            | 3306         | `DB_HOST`, `DB_PORT`, `DB_DATABASE`, `DB_USERNAME`, `DB_PASSWORD`                |
| **PostgreSQL**    | PostgreSQL 15 Database        | 5432         | `DB_HOST`, `DB_PORT`, `DB_DATABASE`, `DB_USERNAME`, `DB_PASSWORD`                |
| **Redis**         | Redis 7 Cache/Session Store   | 6379         | `REDIS_HOST`, `REDIS_PORT`                                                       |
| **MongoDB**       | MongoDB 7 NoSQL Database      | 27017        | `MONGO_HOST`, `MONGO_PORT`, `MONGO_DATABASE`, `MONGO_USERNAME`, `MONGO_PASSWORD` |
| **Kafka**         | Apache Kafka + Zookeeper      | 9092         | `KAFKA_HOST`, `KAFKA_PORT`                                                       |
| **Elasticsearch** | Elasticsearch 8 Search Engine | 9200         | `ELASTICSEARCH_HOST`, `ELASTICSEARCH_PORT`                                       |
| **RabbitMQ**      | RabbitMQ Message Broker       | 5672/15672   | `RABBITMQ_HOST`, `RABBITMQ_PORT`, `RABBITMQ_USERNAME`, `RABBITMQ_PASSWORD`       |

## Available Profiles

| Profile   | Services           | Best For                       |
| --------- | ------------------ | ------------------------------ |
| `laravel` | MySQL + Redis      | Laravel, PHP applications      |
| `node`    | PostgreSQL + Redis | Node.js, Express applications  |
| `django`  | PostgreSQL + Redis | Django, Python applications    |
| `rails`   | PostgreSQL + Redis | Ruby on Rails applications     |
| `spring`  | MySQL + Kafka      | Spring Boot, Java applications |
| `full`    | All services       | Multi-service development      |

## Command Reference

### Core Commands

```bash
dockenv init                    # Interactive setup wizard
dockenv init --profile laravel  # Use predefined profile
dockenv init --services mysql,redis  # Specify services directly

dockenv up                      # Start all services
dockenv up mysql               # Start specific service
dockenv down                   # Stop all services
dockenv restart                # Restart all services
dockenv restart redis          # Restart specific service

dockenv status                 # Show service status
dockenv logs                   # Show all logs
dockenv logs -f mysql          # Follow MySQL logs
```

### Service Management

```bash
dockenv add mysql              # Add MySQL to existing setup
dockenv add redis mongodb      # Add multiple services
dockenv remove postgres        # Remove PostgreSQL
dockenv list                   # Show available services and profiles
```

### Auto-start Management

```bash
dockenv autostart enable       # Enable auto-start on boot
dockenv autostart disable      # Disable auto-start
dockenv autostart status       # Show auto-start status
```

## Configuration

### Directory Structure

```
your-project/
├── docker-compose.dockenv.yaml  # Generated Docker Compose file
├── .env                         # Generated environment variables
└── ~/.config/dockenv/
    └── dockenv.yaml            # dockenv configuration
```

### Configuration File

The configuration is stored in `~/.config/dockenv/dockenv.yaml`:

```yaml
version: "1.0"
services:
  - mysql
  - redis
ports:
  mysql: 3306
  redis: 6379
env:
  DB_CONNECTION: mysql
  DB_HOST: 127.0.0.1
  DB_PORT: "3306"
  DB_DATABASE: dockenv
  DB_USERNAME: dockenv
  DB_PASSWORD: password
  REDIS_HOST: 127.0.0.1
  REDIS_PORT: "6379"
data_path: /home/user/.local/share/dockenv
```

### Custom Ports

```bash
# Set custom ports during initialization
dockenv init --port mysql:3307 --port redis:6380

# Or when adding services
dockenv add --port postgres:5433 postgres
```

### Custom Data Directory

```bash
# Use custom data directory
dockenv init --data-path /path/to/data

# Or set environment variable
export DOCKENV_DATA=/path/to/data
dockenv init
```

## Integration Examples

### Laravel

```bash
# Setup Laravel environment
dockenv init --profile laravel
dockenv up

# Your .env file will contain:
# DB_CONNECTION=mysql
# DB_HOST=127.0.0.1
# DB_PORT=3306
# DB_DATABASE=dockenv
# DB_USERNAME=dockenv
# DB_PASSWORD=password
# REDIS_HOST=127.0.0.1
# REDIS_PORT=6379
```

### Node.js with Express

```bash
# Setup Node.js environment
dockenv init --profile node
dockenv up

# Connect to PostgreSQL
const client = new Client({
  host: process.env.DB_HOST,
  port: process.env.DB_PORT,
  database: process.env.DB_DATABASE,
  user: process.env.DB_USERNAME,
  password: process.env.DB_PASSWORD,
});
```

### Django

```bash
# Setup Django environment
dockenv init --profile django
dockenv up

# settings.py
DATABASES = {
    'default': {
        'ENGINE': 'django.db.backends.postgresql',
        'NAME': os.environ.get('DB_DATABASE'),
        'USER': os.environ.get('DB_USERNAME'),
        'PASSWORD': os.environ.get('DB_PASSWORD'),
        'HOST': os.environ.get('DB_HOST'),
        'PORT': os.environ.get('DB_PORT'),
    }
}
```

## Troubleshooting

### Docker Issues

```bash
# Check Docker installation
dockenv status

# Docker not running
sudo systemctl start docker

# Permission issues (Linux)
sudo usermod -aG docker $USER
# Then logout and login again
```

### Service Connection Issues

```bash
# Check if services are running
dockenv status
docker ps

# Check service logs
dockenv logs mysql
dockenv logs -f redis

# Restart problematic services
dockenv restart mysql
```

### Port Conflicts

```bash
# Use custom ports to avoid conflicts
dockenv add --port mysql:3307 mysql

# Or edit configuration and restart
dockenv restart
```

## Development

### Building from Source

```bash
git clone https://github.com/mohammed-bageri/dockenv.git
cd dockenv
go mod tidy
go build -o dockenv .
```

### Running Tests

```bash
go test ./...
```

### Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Commit changes: `git commit -am 'Add feature'`
4. Push to branch: `git push origin feature-name`
5. Submit a Pull Request

## FAQ

**Q: Can I use dockenv with existing Docker Compose files?**
A: dockenv creates its own compose file (`docker-compose.dockenv.yaml`). You can run both alongside each other with different project names.

**Q: Will dockenv interfere with my existing Docker containers?**
A: No, dockenv uses prefixed container names (`dockenv-mysql`, `dockenv-redis`, etc.) to avoid conflicts.

**Q: Can I customize the Docker images used?**
A: Currently, dockenv uses predefined images optimized for development. Custom image support is planned for future releases.

**Q: How do I backup my data?**
A: Data is stored in Docker volumes. You can backup using:

```bash
docker run --rm -v dockenv_mysql_data:/data -v $(pwd):/backup ubuntu tar czf /backup/mysql-backup.tar.gz /data
```

**Q: Can I use dockenv in production?**
A: No, dockenv is designed for development environments only. Use proper Docker Compose files for production.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Docker](https://docker.com) for containerization
- [Cobra](https://github.com/spf13/cobra) for CLI framework
- [promptui](https://github.com/manifoldco/promptui) for interactive prompts

---

**Made with ❤️ for developers who love simple tools**
