package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func okResponse() []byte {
	return []byte("+OK\r\n")
}

func nullBulkString() []byte {
	return []byte("$-1\r\n")
}

func get(conn net.Conn, key_store KeyStore, get_args []string) {
	fmt.Printf("GET requested for key '%s'\n", get_args[0])
	key_value := key_store[get_args[0]]
	now := time.Now()

	if !key_value.expiration.IsZero() &&
		key_value.expiration.Before(now) {
		fmt.Printf("GET request value at '%s' but expired at '%s'\n",
			now, key_value.expiration)
		conn.Write(nullBulkString())
	} else {
		fmt.Printf("GET request is fresh, expiration at '%s'\n",
			key_value.expiration)
		conn.Write(encodeResponseString(key_value.value))
	}
}

func set(conn net.Conn, key_store KeyStore, set_args []string) {
	fmt.Printf("SET requested: key '%s', value '%s'\n",
		set_args[0], set_args[1])

	now := time.Now()
	key_value := KeyValue{
		value:      set_args[1],
		expiration: time.Time{},
	}

	if len(set_args) >= 3 {

		arg := strings.ToLower(strings.TrimSpace(set_args[2]))
		if arg == "px" {
			ms, err := strconv.Atoi(strings.TrimSpace(set_args[3]))

			if err == nil {
				key_value.expiration = now.Add(
					time.Second * time.Duration(ms) / 1000)
				fmt.Printf(
					"SET expiration for key '%s' at +%dms to '%s'\n",
					key_value.value, ms, key_value.expiration)
			} else {
				fmt.Printf("Error parsing PX expiry value '%s'\n", arg)
			}
		}
	}

	key_store[set_args[0]] = key_value
	conn.Write(okResponse())
}

func echo(conn net.Conn, echo_args []string) {
	fmt.Println("ECHO requested")
	response := encodeResponseString(
		strings.Join(echo_args, " "))
	fmt.Printf("Sending back: '%s", response)
	conn.Write(response)
}

func encodeResponseString(response string) []byte {
	length := len(response)

	terminator := "\r\n"
	response_bytes := make([]byte, 0)

	response_bytes = append(response_bytes, '$')
	response_bytes = append(response_bytes, []byte(strconv.Itoa(length))...)

	response_bytes = append(response_bytes, []byte(terminator)...)
	response_bytes = append(response_bytes, []byte(response)...)

	response_bytes = append(response_bytes, []byte(terminator)...)

	return response_bytes
}

func ping(conn net.Conn) {
	fmt.Println("PING requested")
	fmt.Println("writing PONG response")
	conn.Write([]byte("+PONG\r\n"))
}
