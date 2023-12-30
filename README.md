Overview

This is a small server managing a twitter like an app

API Documentation
Create User (POST /users)
Create a new user.
To create a new user copy this bash command, feel free to choose any id and any username:

curl -X POST -H "Content-Type: application/json" -d '{"id": "1", "username": "user1", "following": []}' http://localhost:8080/users

Find User (GET /users?id={userID})
Retrieve user details by ID.

curl http://localhost:8080/users?id='userId'

Update User (PUT /users?id={userID})
Update user details by ID.

curl -X PUT -H "Content-Type: application/json" -d '{"id": "1", "username": "updatedUser", "following": []}' http://localhost:8080/users?id=1

Post Tweet (POST /tweets)
Post a new tweet.

curl -X POST -H "Content-Type: application/json" -d '{"id": "1", "userId": "user1", "content": "Hello, Twitter!", "timestamp": "2023-01-01T12:00:00Z"}' http://localhost:8080/tweets

Follow User (POST /follow)
Make a user follow another user.

curl -X POST http://localhost:8080/follow?follower_id=1&followed_id=2

Unfollow User (DELETE /follow)
Make a user unfollow another user.

curl -X DELETE http://localhost:8080/follow?follower_id=1&unfollowed_id=2

User Wall (GET /wall?id={userID})
Retrieve a user's wall.

curl http://localhost:8080/wall?id=1

