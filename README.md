# Go Redis
** A redis implementation using golang. **
Go redis is an implementation of Redis using golang from scratch. 
The project supports handling of GET, SET commands using a TCP connection like redis.

#### This project familiarized me to the "net" package and using "tcp" connections. Before this I'd only been using the "net/http" package to spin up a webserver. 
But implementing and handling the connections at the "tcp" level allowed me to understand Go more deeply.

It also includes tests showcasing clients setting, getting the key values.

## Installation
```bash
# 1. Clone the repository
git clone https://github.com/manavkush/Go-redis.git
cd Go-redis
# 2. Get the dependencies
go mod tidy
```

## Running the server
```
make run
```

## Changing listening ports
You can open the Makefile and change the port to specify the port that you want the redis server to listen to requests.


