# redis-pipe [![Build Status](https://travis-ci.org/lukasmartinelli/redis-pipe.svg?branch=master)](https://travis-ci.org/lukasmartinelli/redis-pipe)

**redis-pipe** allows you to treat [Redis Lists](http://redis.io/topics/data-types#lists)
as if they were [Unix pipes](https://en.wikipedia.org/wiki/Pipeline_%28Unix%29).
It basically connects `stdin` and `stdout` with `LPUSH` and `LPOP`.

## How it works

### Configuration

Set the `REDIS_HOST` and `REDIS_PORT` environment variables for
easy configuration or pass `--host` and `--port` arguments.

### Writing from stdin to Redis List

Pipe in value to `redis-pipe` and it will `LPUSH` them to the Redis List.

```
echo "hi there" | ./redis-pipe greetings
```

![Write from stdin to Redis with LPUSH](redis-lpush.png)

### Reading from Redis List to stdout

If you call `redis-pipe` with a tty attached it will `LPOP` all values
from the Redis List and write them to stdout.

```
./redis-pipe greetings
```

![Read from Redis with LPOP and write to stdout](redis-lpop.png)

You can also limit the amount of values popped from the list.

```
./redis-pipe --count 100 greetings
```

Support for blocking mode with `BLPOP` is not supported yet.

## Examples

### Centralized Logging

In this sample we pipe the syslog to a Redis List called `log`.

```
tail -f /var/log/syslog | ./redis-pipe logs
```

You can now easily collect all the syslogs of your machines
on a single server.

```
./redis-pipe logs > logs.txt
```

### Very basic job queue

Create jobs and store them.

```
cat jobs.txt | redis-pipe jobs
```

Process jobs on several workers and store the results.

```
./redis-pipe --count 10 jobs | python do-work.py | ./redis-pipe results
```

Collect the results.
```
./redis-pipe results > results.txt
```

## Install

Simply download the release and extract it.

### OSX

```
wget https://github.com/lukasmartinelli/redis-pipe/releases/download/v1.4.1/redis-pipe_darwin_amd64.zip
unzip redis-pipe_darwin_amd64.zip
cd redis-pipe_darwin_amd64
./redis-pipe --help
```

### Linux

```
wget https://github.com/lukasmartinelli/redis-pipe/releases/download/v1.4.1/redis-pipe_linux_amd64.tar.gz
tar -xvzf redis-pipe_linux_amd64.tar
cd redis-pipe_linux_amd64
./redis-pipe --help
```

## Build

Install dependencies

```
go get github.com/andrew-d/go-termutil
go get github.com/codegangsta/cli
go get github.com/garyburd/redigo/redis
```

Build binary

```
go build redis-pipe.go
```
