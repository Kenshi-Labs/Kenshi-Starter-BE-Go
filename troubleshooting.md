# Troubleshooting Guide

## MongoDB Connection Issues
1. Check if MongoDB is running:
```bash
docker ps | grep auth-mongodb
```

2. Check MongoDB logs:
```bash
docker logs auth-mongodb
```

3. Verify MongoDB connection:
```bash
docker exec -it auth-mongodb mongosh
```

## API Issues
1. Check API logs:
```bash
# Run with detailed logging
go run main.go 2>&1 | tee api.log
```

2. Common HTTP Status Codes:
- 400: Bad Request (invalid input)
- 401: Unauthorized (invalid/missing token)
- 403: Forbidden (insufficient permissions)
- 404: Not Found
- 409: Conflict (duplicate email)
- 500: Internal Server Error

3. JWT Token Issues:
- Make sure token is properly formatted in Authorization header
- Check token expiration
- Verify JWT_SECRET is consistent

## Database Issues
1. Check collections:
```javascript
use auth_db
db.users.find()  // List users
db.roles.find()  // List roles
```

2. Check indexes:
```javascript
db.users.getIndexes()
```

3. Reset database if needed:
```javascript
use auth_db
db.users.drop()
db.roles.drop()
```