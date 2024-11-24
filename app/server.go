package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/pborman/getopt/v2"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

type KeyValue struct {
	value      string
	expiration time.Time
}

type Config map[string]string

type KeyStore map[string]KeyValue

func config_string(config Config) string {
	var str strings.Builder

	for k, v := range config {
		s := fmt.Sprintf("%-20s%s\n", k, v)
		str.WriteString(s)
	}
	return str.String()
}

func main() {

	var dir = "/tmp/redis-data"
	var dbfilename = "dump.rdb"

	config := Config{}
	getopt.FlagLong(&dir,
		"dir", 0, "Redis database directory")
	getopt.FlagLong(&dbfilename,
		"dbfilename", 0, "Redis database filename")

	getopt.Parse()
	key_store := KeyStore{}
	fmt.Println(dir)
	config["dir"] = dir
	config["dbfilename"] = dbfilename

	fmt.Println("SuperRedis is ready to GO!")
	fmt.Println("Config:")
	fmt.Println(config_string(config))

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		fmt.Println("Connection established")

		go handle(conn, config, key_store)

	}

}
