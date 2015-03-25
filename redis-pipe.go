package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/andrew-d/go-termutil"
	"github.com/docopt/docopt-go"
	"menteslibres.net/gosexy/redis"
)

const usage string = `
Treat Redis Lists like Unix Pipes by connecting stdin with LPUSH
and LPOP with stout.

Usage:
    redis-pipe <list> [-n <count>] [--blocking] [--host=<host>] [--port=<port>]
    redis-pipe (-h | --help)

Options:
    -h --help         Show this screen
    --blocking        Read in blocking mode from list (infinite timeout)
	--count=<count>   Stop reading from list after count [default: -1]
    --host=<host>     Redis host [default: localhost]
    --port=<port>     Redis port [default: 6379]
`

//Pop all values from the given redis list and write the values to stdout
func readAll(list string, client *redis.Client, blocking bool) {
	for {
		read(list, client, 1, blocking)
	}
}

func read(list string, client *redis.Client, count int, blocking bool) {
	for i := 0; i < count; i++ {
		var values []string
		var err error

		if blocking {
			values, err = client.BLPop(0, list)
		} else {
			values, err = client.BLPop(1, list)
		}

		if err != nil {
			os.Exit(0)
		} else {
			fmt.Println(values[1])
		}
	}
}

//Push all values from stdin to a given redis list
func write(list string, client *redis.Client) {
	scanner := bufio.NewScanner(os.Stdin)
	values := make([]interface{}, 64)
	for i := 0; i < 64; i++ {
		for scanner.Scan() {
			values = append(values, scanner.Text())
		}
	}
	_, err := client.LPush(list, values...)
	if err != nil {
		log.Fatalf("Could not LPUSH: %s \"%s\"\n", list, err.Error())
	}
}

//Return host and port of redis server inferred from either
//environment variables or command line argumentsof redis server
func redisConfig(args map[string]interface{}) (string, uint) {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	if host == "" {
		host = args["--host"].(string)
	}

	if port == "" {
		port = args["--port"].(string)
	}

	portNum, _ := strconv.ParseUint(port, 10, 6379)
	return host, uint(portNum)
}

func main() {
	args, _ := docopt.Parse(usage, nil, true, "pusred 0.9", false)
	list := args["<list>"].(string)
	blockingMode := args["--blocking"].(bool)
	host, port := redisConfig(args)
	countStr, _ := args["--count"].(string)
	count, _ := strconv.ParseInt(countStr, 10, -1)

	client := redis.New()
	err := client.Connect(host, port)
	if err != nil {
		log.Fatalf("Connecting to Redis failed: %s\n", err.Error())
	}

	if termutil.Isatty(os.Stdin.Fd()) {
		if count > 0 {
			read(list, client, int(count), blockingMode)
		} else {
			readAll(list, client, blockingMode)
		}
	} else {
		write(list, client)
	}

	client.Quit()
}
