# redis-streams-with-go
This is a tutorial on how to effectively reading Redis streams requires some work: counting ids, 
prefetching and buffering, asynchronously sending acknowledgements and parsing entries. 

## Tools to be used 
* 
* [Redis](https://redis.io)
* [Redis Streams](https://redis.io/docs/latest/develop/data-types/streams/)
* [Golang](https://go.dev)
* [Golang Redis client](https://redis.uptrace.dev)
* [Go Typed Redis Streams](https://github.com/dranikpg/gtrs)
* [Docker](https://www.docker.com)


## How to build
You can use your favorite editor to code (e.g. IntelliJ, Visual Code, Emacs, etc)

To build the project you need docker-compose and docker engine installed in you pc

```shell
# do build the code 
docker-compose build 

# to run the producers, consumers, redis
docker-compose up

# to shutdown the docker images
docker-compose down
```

## Redis Streams Cheat Sheet

### Adding data to the stream 

## Consuming data from the stream

## Consuming data from the stream via consumer groups

## Resources

