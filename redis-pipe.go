package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/andrew-d/go-termutil"
	"github.com/codegangsta/cli"
	"menteslibres.net/gosexy/redis"
)

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
			// The 1st value is the list name and only the 2nd is the actual value
			value := values[1]
			value = strings.TrimSpace(value)
			if value != "" {
				fmt.Println(value)
			}
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
	app := cli.NewApp()
	app.Name = "redis-pipe"
	app.Usage = "connect stdin with LPUSH and LPOP with stout."
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "host",
			Value:  "localhost",
			Usage:  "Redis host",
			EnvVar: "REDIS_HOST",
		},
		cli.IntFlag{
			Name:   "port",
			Value:  6379,
			Usage:  "Redis port",
			EnvVar: "REDIS_PORT",
		},
		cli.IntFlag{
			Name:  "count",
			Value: -1,
			Usage: "Redis Stop reading from list after count",
		},
		cli.BoolFlag{
			Name:  "blocking",
			Usage: "Read in blocking mode from list (infinite timeout)",
		},
	}

	app.Action = func(c *cli.Context) {
		list := c.Args().First()
		if list == "" {
			fmt.Println("Please provide name of the list")
			os.Exit(1)
		}

		client := redis.New()
		err := client.Connect(c.String("host"), uint(c.Int("port")))
		if err != nil {
			log.Fatalf("Connecting to Redis failed: %s\n", err.Error())
		}

		if termutil.Isatty(os.Stdin.Fd()) {
			count := c.Int("count")
			blockingMode := c.Bool("blocking")
			if count > 0 {
				read(list, client, count, blockingMode)
			} else {
				readAll(list, client, blockingMode)
			}
		} else {
			write(list, client)
		}
		client.Quit()
	}

	app.Run(os.Args)
}
