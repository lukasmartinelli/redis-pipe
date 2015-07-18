# redis-pipe [![Build Status](https://travis-ci.org/lukasmartinelli/redis-pipe.svg?branch=master)](https://travis-ci.org/lukasmartinelli/redis-pipe)

**redis-pipe** allows you to treat [Redis Lists](http://redis.io/topics/data-types#lists)
as if they were [Unix pipes](https://en.wikipedia.org/wiki/Pipeline_%28Unix%29).
It basically connects `stdin` and `stdout` with `LPUSH` and `LPOP`.

## Installing from source

```shell
$ go get github.com/lukasmartinelli/redis-pipe
```

* To be able to use the the built binary in any shell,

make sure to have your $GOPATH properly set in your

   ~/.bash_profile (on OS X)

or

   ~/.bash_rc (on Linux)

file e.g:

```shell
$ cat << ! >> ~/.bash_rc
> export GOPATH="\$HOME/gopath"
> export PATH="\$GOPATH:\$GOPATH/bin:\$PATH"
> !
$ source ~/.bash_rc
```

## How it works

### Configuration

Set the `REDIS_HOST` and `REDIS_PORT` environment variables for
easy configuration or pass `--host` and `--port` arguments.

### Writing from stdin to Redis List

Pipe in value to `redis-pipe` and it will `LPUSH` them to the Redis List.

```
echo "hi there" | redis-pipe greetings
```

![Write from stdin to Redis with LPUSH](redis-lpush.png)

### Reading from Redis List to stdout

If you call `redis-pipe` with a tty attached it will `LPOP` all values
from the Redis List and write them to stdout.

```
redis-pipe greetings
```

![Read from Redis with LPOP and write to stdout](redis-lpop.png)

You can also limit the amount of values popped from the list.

```
redis-pipe --count 100 greetings
```

Support for blocking mode with `BLPOP` is not supported yet.

## Examples

### Centralized Logging

In this sample we pipe the syslog to a Redis List called `logs`.

```
tail -f /var/log/syslog | redis-pipe logs
```

You can now easily collect all the syslogs of your machines
on a single server.

```
redis-pipe logs > logs.txt
```

### Very basic job queue

Create jobs and store them.

```
cat jobs.txt | redis-pipe jobs
```

Process jobs on several workers and store the results.

```
redis-pipe --count 10 jobs | python do-work.py | redis-pipe results
```

Collect the results.
```
redis-pipe results > results.txt
```

## Installing from binary releases

Simply download the release and extract it.
Add the binary install path to your ~/.bash_rc or ~/.bash_profile file
e.g

### OSX

```
wget https://github.com/lukasmartinelli/redis-pipe/releases/download/v1.4.1/redis-pipe_darwin_amd64.zip
unzip redis-pipe_darwin_amd64.zip
cd redis-pipe_darwin_amd64
cat << $ >> ~/.bash_profile
> export PATH = "$(pwd):\$PATH"
> $
$ source ~/.bash_profile
redis-pipe --help
```

### Linux

```
wget https://github.com/lukasmartinelli/redis-pipe/releases/download/v1.4.1/redis-pipe_linux_amd64.tar.gz
tar -xvzf redis-pipe_linux_amd64.tar
cd redis-pipe_linux_amd64
cat << $ >> ~/.bash_rc
> export PATH = "$(pwd):\$PATH"
> $
$ source ~/.bash_rc
redis-pipe --help
```
