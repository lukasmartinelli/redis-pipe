# redis-pipe [![Build Status](https://travis-ci.org/lukasmartinelli/redis-pipe.svg?branch=master)](https://travis-ci.org/lukasmartinelli/redis-pipe)

**redis-pipe** allows you to treat [Redis Lists](http://redis.io/topics/data-types#lists)
as if they were [Unix pipes](https://en.wikipedia.org/wiki/Pipeline_%28Unix%29).
It basically connects `stdin` and `stdout` with `LPUSH` and `LPOP`.

## Build

Install dependencies

```
go get github.com/andrew-d/go-termutil
go get github.com/docopt/docopt-go
go get menteslibres.net/gosexy/redis
```

Build binary

```
go build redis-pipe.go
```

## How it works

### Writing from stdin to Redis List

**redis-pipe** takes your values and generates `RPUSH` commands
(generating a valid [Redis protocol](http://redis.io/topics/protocol))
that are then piped into `redis-cli --pipe` ([Redis Mass Insertion](http://redis.io/topics/mass-insert))

### Reading from Redis List to stdout

`LPOP` all the values from the list and write them to `stdout`.

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

## Usage

You can set the `REDIS_HOST` and `REDIS_PORT` environment variables for
easy configuration.
