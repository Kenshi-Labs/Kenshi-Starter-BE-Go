# Kenshi Starter BE API (Go)

A production-ready authentication API built with Go Fiber, MongoDB, and JWT. Features include user authentication, RBAC (Role-Based Access Control), refresh tokens, and password reset functionality. This will be converted into a full blown starter kit as mentioned in the starter_kit.md file

## Features

- 🔐 JWT Authentication
- 🔄 Refresh Tokens
- 👥 Role-Based Access Control (RBAC)
- 🔑 Password Reset with Email
- 🚀 Live Reload Development Environment
- 🐳 Docker-based Development Setup
- 📧 Email Testing with MailHog
- 📊 MongoDB Admin Interface

## Prerequisites

- Docker and Docker Compose
- Git
- Make (optional, for using Makefile commands)

## Project Structure

```
├── configs/         # Configuration files
├── handlers/        # Request handlers
├── middleware/      # Middleware functions
├── models/         # Data models
├── utils/          # Utility functions
├── .air.toml       # Air configuration for live reload
├── .env            # Environment variables
├── docker-compose.yml
├── Dockerfile.dev   # Development Dockerfile
├── go.mod          # Go modules file
├── main.go         # Application entry point
└── README.md       # This file
```

## Quick Start

1. Clone the repository:
```bash
git clone https://github.com/yourusername/go-auth-api.git
cd go-auth-api
```

2. Create `.env` file:
```bash
cp .env.example .env
```

3. Start the development environment:
```bash
# Development environment
docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build

# Production with monitoring
docker compose -f docker-compose.yml -f docker-compose.prod.yml -f docker-compose.monitoring.yml up -d  --build

# To down specific compose file
docker compose -f docker-compose.yml -f docker-compose.monitoring.yml down

# To view logs from all services
docker compose -f docker-compose.yml -f docker-compose.monitoring.yml logs -f

# To rebuild specific services
docker compose -f docker-compose.yml -f docker-compose.dev.yml build api
```

The application will be available at:
- API: http://localhost:3000
- MongoDB Express: http://localhost:8081
- MailHog (Email Testing): http://localhost:8025

## Available Endpoints

### Authentication
```
POST /api/auth/signup
POST /api/auth/signin
POST /api/auth/refresh
```

### Password Reset
```
POST /api/auth/forgot-password
POST /api/auth/reset-password
```

### User Management (Protected Routes)
```
GET  /api/user/profile
PUT  /api/user/profile
DELETE /api/user/profile
```

## Development

### Live Reload
The development environment uses Air for live reloading. Any changes to your Go files will automatically rebuild and restart the application.

### Database Management
- MongoDB Express is available at http://localhost:8081
- Initial database setup is handled by init-mongo.js
- Database persists in a Docker volume

### Email Testing
All emails in development are caught by MailHog and can be viewed at http://localhost:8025

## Environment Variables

```
MONGODB_URI=mongodb://mongo:27017/auth_db
JWT_SECRET=your-dev-secret-key
REFRESH_SECRET=your-refresh-secret-key
API_PORT=3000
SMTP_HOST=mailhog
SMTP_PORT=1025
SMTP_USER=
SMTP_PASS=
APP_URL=http://localhost:3000
```

## Testing API Endpoints

### Sign Up
```bash
curl -X POST http://localhost:3000/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### Sign In
```bash
curl -X POST http://localhost:3000/api/auth/signin \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### Using Protected Routes
```bash
# Get Profile (use token from signin response)
curl -X GET http://localhost:3000/api/user/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Common Issues

### MongoDB Connection Issues
If you can't connect to MongoDB, ensure:
1. MongoDB container is running: `docker ps`
2. MongoDB URI is correct in .env
3. Try restarting the containers: `docker compose down && docker compose up`

### Live Reload Not Working
1. Check air.log file in the project root
2. Ensure your .air.toml configuration is correct
3. Try rebuilding: `docker compose build --no-cache`

## Contributing

1. Fork the repository
2. Create your feature branch: `git checkout -b feature/my-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin feature/my-feature`
5. Submit a pull request

## Directory Structure Explanation

- `configs/`: Configuration management and database connection
- `handlers/`: HTTP request handlers for each route
- `middleware/`: Custom middleware (auth, RBAC, validation)
- `models/`: Data structures and database models
- `utils/`: Helper functions and utilities

## Security Notes

- JWT tokens expire after 15 minutes
- Refresh tokens are valid for 7 days
- Passwords are hashed using bcrypt
- All routes are HTTPS-only in production
- CORS is enabled for specified origins only

## License

MIT License - see LICENSE file for details