

This is a simple Http server which is protected by authentication.
- Use Json file to simulate the database behaviors
- Use go-redis to implement the redis layer
- User's password is encrypted in http request and stored in database 


## Deploy Redis with docker
docker run -itd --name redis-test -p 6379:6379 redis

## How to run this server
- git clone https://github.com/sunportainer/air
- cd ./air
- go run main.go

## Send request with curl
- curl http://{IPAddress}:3333?email=email&pwd=encryptedPassword to get response

- curl -d 'firstName={fName}&lastName={lName}&email={email}&pwd={passwd}' http://{IPAddress}:3333 to create new user

## Note
- {passwd} string should be encrypted with any available hash method and be URL encoded.
