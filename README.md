# pusred

A simple Redis client that connects stdin and stdout with `LPUSH` and `LPOP`.

## Examples

Add value to list

```
echo "Hello" | ./pusred lpush greetings
```

Pipe syslog to Redis

```
tail -f /var/log/syslog | ./pusred lpush logs
```

Read all values from list

```
./pusred lpop greetings
```

## Usage

```
Usage:
    pusred lpop <list> [--host=<host>] [--port=<port>]
    pusred lpush <list> [--host=<host>] [--port=<port>]
    pusred (-h | --help)

Options:
    -h --help         Show this screen
    --host=<host>     Redis host [default: localhost]
    --port=<port>     Redis port [default: 6379]
```

You can also set the environment variables `RQ_HOST` and `RQ_PORT` to
define the Redis connection.
