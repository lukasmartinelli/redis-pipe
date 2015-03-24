package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/docopt/docopt-go"
	"menteslibres.net/gosexy/redis"
)

const usage string = `
Push and pop values from Redis lists.

Usage:
    pusred lpop <list> [--host=<host>] [--port=<port>]
    pusred lpush <list> [--host=<host>] [--port=<port>]
    pusred (-h | --help)

Options:
    -h --help         Show this screen
    --host=<host>     Redis host [default: localhost]
    --port=<port>     Redis port [default: 6379]
`

//Pop all values from the given redis list and write the values to stdout
func lpop(list string, client *redis.Client) {
	length, _ := client.LLen(list)
	for length > 0 {
		value, err := client.LPop(list)
		if err != nil {
			log.Fatalf("Could not LPOP%s: %s\n", list, err.Error())
		} else {
			fmt.Println(value)
		}
		length, _ = client.LLen(list)
	}
}

//Push all values from stdin to a given redis list
func lpush(list string, client *redis.Client) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		value := scanner.Text()
		_, err := client.LPush(list, value)
		if err != nil {
			log.Fatalf("Could not LPUSH: %s \"%s\"\n", list, err.Error())
		}
	}
}

//Return host and port of redis server inferred from either
//environment variables or command line argumentsof redis server
func redisConfig(args map[string]interface{}) (string, uint) {
	host := os.Getenv("RQ_HOST")
	port := os.Getenv("RQ_PORT")

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
	host, port := redisConfig(args)

	client := redis.New()
	err := client.Connect(host, port)
	if err != nil {
		log.Fatalf("Connecting to Redis failed: %s\n", err.Error())
	}

	if args["lpush"].(bool) {
		lpush(list, client)
	}

	if args["lpop"].(bool) {
		lpop(list, client)
	}

	client.Quit()
}
