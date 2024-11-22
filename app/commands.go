package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func okResponse() []byte {
	return []byte("+OK\r\n")
}

func get(conn net.Conn, key_store KeyStore, get_args []string) {
	fmt.Printf("GET requested for key '%s'", get_args[0])
	response := key_store[get_args[0]]
	conn.Write(encodeResponseString(response))
}

func set(conn net.Conn, key_store KeyStore, set_args []string) {
	fmt.Printf("SET requested: key '%s', value '%s'",
		set_args[0], set_args[1])
	key_store[set_args[0]] = set_args[1]
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
