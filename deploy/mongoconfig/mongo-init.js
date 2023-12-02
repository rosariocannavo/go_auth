db = db.getSiblingDB('my_database');


db.createUser(
    {
        user: "myuser",
        pwd: "mypassword",
        roles: [
            {
                role: "readWrite",
                db: "my_database"
            }
        ]
    }
);

db.createCollection('users');




