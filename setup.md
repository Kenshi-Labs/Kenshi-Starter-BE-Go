# Development and Production Setup Guide

## Development Setup

1. Clone the repository and create environment file:
```bash
# Clone repository
git clone 
cd 

# Create .env file for development
cat > .env << EOL
MONGODB_URI=mongodb://mongo:27017/auth_db
JWT_SECRET=dev-secret-key
REFRESH_SECRET=dev-refresh-key
API_PORT=3000
SMTP_HOST=mailhog
SMTP_PORT=1025
EOL
```

2. Start development environment:
```bash
# Start all services with hot reload
docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d

# View logs
docker compose logs -f api

# Access different services:
# - API: http://localhost:3000
# - MongoDB Express: http://localhost:8081
# - MailHog: http://localhost:8025
```

3. Development Commands:
```bash
# Rebuild specific service
docker compose -f docker-compose.yml -f docker-compose.dev.yml build api

# Restart service
docker compose -f docker-compose.yml -f docker-compose.dev.yml restart api

# View logs of specific service
docker compose logs -f api

# Stop all services
docker compose -f docker-compose.yml -f docker-compose.dev.yml down
```

## Production Deployment

1. Server Prerequisites:
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install required packages
sudo apt install -y docker.io docker-compose nginx certbot

# Start and enable Docker
sudo systemctl start docker
sudo systemctl enable docker
```

2. Setup SSL Certificate:
```bash
# Stop nginx if running
sudo systemctl stop nginx

# Get SSL certificate
sudo certbot certonly --standalone -d your-domain.com

# Create SSL directory
sudo mkdir -p /etc/nginx/ssl
sudo cp /etc/letsencrypt/live/your-domain.com/fullchain.pem /etc/nginx/ssl/cert.pem
sudo cp /etc/letsencrypt/live/your-domain.com/privkey.pem /etc/nginx/ssl/key.pem
```

3. Create Production Environment File:
```bash
# Create .env.prod file
cat > .env.prod << EOL
MONGODB_URI=mongodb://your-production-mongodb-uri
JWT_SECRET=your-secure-production-secret
REFRESH_SECRET=your-secure-refresh-secret
API_PORT=3000
SMTP_HOST=your-smtp-host
SMTP_PORT=587
SMTP_USER=your-smtp-user
SMTP_PASS=your-smtp-password
GRAFANA_PASSWORD=your-secure-grafana-password
EOL
```

4. Deploy Application:
```bash
# Copy files to production server
scp -r ./* user@your-server:/path/to/app
scp .env.prod user@your-server:/path/to/app/.env

# SSH into server
ssh user@your-server

# Navigate to app directory
cd /path/to/app

# Start production services
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Start monitoring (optional)
docker compose -f docker-compose.monitoring.yml up -d
```

5. Monitoring Setup:
```bash
# Access monitoring tools:
# - Prometheus: http://your-domain:9090
# - Grafana: http://your-domain:3001
# - Node Exporter metrics: http://your-domain:9100/metrics

# Initial Grafana login:
# Username: admin
# Password: [value from GRAFANA_PASSWORD in .env]
```

## Maintenance Commands

1. Update Application:
```bash
# Pull latest changes
git pull origin main

# Rebuild and restart services
docker compose -f docker-compose.yml -f docker-compose.prod.yml build
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

2. View Logs:
```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f api
```

3. Backup Database:
```bash
# Create backup
docker compose exec mongo mongodump --out=/data/db/backup

# Copy backup from container
docker cp $(docker compose ps -q mongo):/data/db/backup ./backup
```

4. Restore Database:
```bash
# Copy backup to container
docker cp ./backup $(docker compose ps -q mongo):/data/db/

# Restore backup
docker compose exec mongo mongorestore /data/db/backup
```

## Health Checks

1. API Health:
```bash
curl http://localhost:3000/health

# Production
curl https://your-domain.com/health
```

2. Monitor Container Status:
```bash
docker compose ps
```

3. Check Resource Usage:
```bash
docker stats
```

## Troubleshooting

1. Container Issues:
```bash
# Check container logs
docker compose logs -f [service-name]

# Restart specific service
docker compose restart [service-name]

# Rebuild and restart service
docker compose up -d --build [service-name]
```

2. Database Issues:
```bash
# Check MongoDB logs
docker compose logs mongo

# Connect to MongoDB shell
docker compose exec mongo mongosh
```

3. Network Issues:
```bash
# Check networks
docker network ls

# Inspect network
docker network inspect [network-name]
```

## Scaling (Production)

```bash
# Scale API service
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --scale api=3

# Check running instances
docker compose ps
```
