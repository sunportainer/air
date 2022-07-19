

This is a simple Http server which is protected by authentication.
- Use Json file to simulate the database behaviors
- Use go-redis to implement the redis layer
- User's password is encrypted in http request and stored in database 


## Deploy Redis with docker
docker run -itd --name redis-test -p 6379:6379 redis
