# Example: Laravel Project Setup

This example shows how to set up a complete Laravel development environment with dockenv.

## Prerequisites

- Docker and Docker Compose installed
- dockenv installed (`curl -s https://raw.githubusercontent.com/mohammed-bageri/dockenv/main/install.sh | bash`)

## Setup

1. **Navigate to your Laravel project directory:**

   ```bash
   cd /path/to/your/laravel-project
   ```

2. **Initialize dockenv with Laravel profile:**

   ```bash
   dockenv init --profile laravel
   ```

   This will:

   - Configure MySQL (port 3306)
   - Configure Redis (port 6379)
   - Generate `docker-compose.dockenv.yaml`
   - Create `.env` with database credentials
   - Set up persistent data volumes

3. **Start the services:**

   ```bash
   dockenv up
   ```

4. **Update your Laravel `.env` file:**

   ```env
   DB_CONNECTION=mysql
   DB_HOST=127.0.0.1
   DB_PORT=3306
   DB_DATABASE=dockenv
   DB_USERNAME=dockenv
   DB_PASSWORD=password

   REDIS_HOST=127.0.0.1
   REDIS_PORT=6379
   REDIS_PASSWORD=
   ```

5. **Run Laravel migrations:**
   ```bash
   php artisan migrate
   ```

## Daily Workflow

```bash
# Start your development day
dockenv up

# Check service status
dockenv status

# View logs if needed
dockenv logs
dockenv logs -f mysql  # Follow MySQL logs

# End your development day
dockenv down
```

## Adding More Services

```bash
# Add Elasticsearch for search functionality
dockenv add elasticsearch

# Add RabbitMQ for queues
dockenv add rabbitmq

# Remove a service you no longer need
dockenv remove mongodb
```

## Auto-start on Boot

```bash
# Enable auto-start (requires sudo)
dockenv autostart enable

# Check status
dockenv autostart status

# Disable if needed
dockenv autostart disable
```

## Connection Strings

After running `dockenv up`, you can connect to services using:

- **MySQL**: `mysql://dockenv:password@localhost:3306/dockenv`
- **Redis**: `redis://localhost:6379`

## Troubleshooting

### Port Conflicts

If default ports are in use, customize them:

```bash
dockenv add --port mysql:3307 mysql
```

### Service Won't Start

Check logs and status:

```bash
dockenv status
dockenv logs mysql
```

### Reset Everything

To start fresh:

```bash
dockenv down --volumes  # WARNING: Deletes all data!
dockenv init --profile laravel
dockenv up
```
