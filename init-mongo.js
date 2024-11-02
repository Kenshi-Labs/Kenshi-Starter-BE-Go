// init-mongo.js
db = db.getSiblingDB('auth_db');

// Create roles collection and add default roles
db.roles.insertMany([
    {
        name: "admin",
        permissions: [
            "read:profile",
            "update:profile",
            "delete:profile",
            "create:user",
            "list:users"
        ]
    },
    {
        name: "user",
        permissions: [
            "read:profile",
            "update:profile"
        ]
    }
]);

// Create unique index on email field
db.users.createIndex({ "email": 1 }, { unique: true });
db.refresh_tokens.createIndex({ "token": 1 }, { unique: true });
db.refresh_tokens.createIndex({ "expires_at": 1 }, { expireAfterSeconds: 0 });
db.refresh_tokens.createIndex({ "user_id": 1 });