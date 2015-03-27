package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/andrew-d/go-termutil"
	"github.com/codegangsta/cli"
	"github.com/garyburd/redigo/redis"
)

//Pop all values from the given redis list and write the values to stdout
func readAll(list string, client redis.Conn) {
	for {
		read(list, client, 64)
	}
}

func read(list string, conn redis.Conn, count int) {
	for i := 0; i < count; i++ {
		conn.Send("LPOP", list)
	}
	conn.Flush()
	for i := 0; i < count; i++ {
		value, err := redis.String(conn.Receive())

		if err != nil {
			os.Exit(0)
		}
		if value != "" {
			fmt.Println(value)
		}
	}
}

func writeAll(list string, client redis.Conn) {
	for {
		write(list, client, 16)
	}
}

//Push all values from stdin to a given redis list
func write(list string, conn redis.Conn, count int) {
	scanner := bufio.NewScanner(os.Stdin)
	for i := 0; i < count; i++ {
		for scanner.Scan() {
			value := strings.TrimSpace(scanner.Text())
			if value != "" {
				conn.Send("LPUSH", list, value)
			}
		}
	}
	conn.Flush()
	for i := 0; i < count; i++ {
		_, err := redis.Int(conn.Receive())
		if err != nil {
			log.Fatalf("Could not LPUSH: %s \"%s\"\n", list, err.Error())
		}
	}
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
	}

	app.Action = func(c *cli.Context) {
		list := c.Args().First()
		if list == "" {
			fmt.Println("Please provide name of the list")
			os.Exit(1)
		}

		conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", c.String("host"), c.Int("port")))
		if err != nil {
			log.Fatalf("Connecting to Redis failed: %s\n", err.Error())
		}

		if termutil.Isatty(os.Stdin.Fd()) {
			count := c.Int("count")
			if count > 0 {
				read(list, conn, count)
			} else {
				readAll(list, conn)
			}
		} else {
			writeAll(list, conn)
		}
		conn.Close()
	}

	app.Run(os.Args)
}
