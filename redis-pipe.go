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

func writeBatch(list string, conn redis.Conn, values []string) {
	conn.Send("MULTI")
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			conn.Send("LPUSH", list, value)
		}
	}
	_, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		log.Fatalf("Could not LPUSH: %s \"%s\"\n", list, err.Error())
	}
}

//Push all values from stdin to a given redis list
func write(list string, conn redis.Conn) {
	var values []string
	count := 0
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		value := scanner.Text()
		values = append(values, value)
		count++

		if count >= 64 {
			writeBatch(list, conn, values)
			values = make([]string, 0)
			count = 0
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Reading standard input:", err)
	}
	//writeBatch(list, conn, values)
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
			write(list, conn)
		}
		conn.Close()
	}

	app.Run(os.Args)
}
