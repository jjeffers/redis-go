package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

func handle(conn net.Conn, config Config, key_store KeyStore) {
	defer conn.Close()

	fmt.Println("Handling connection")

	for {

		tmp := make([]byte, 1)

		_, err := conn.Read(tmp)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client closed the connection")

			}
			fmt.Printf("Error reading first byte of request %s\n", err)
		}

		switch tmp[0] {
		case '*':
			bulkArray(conn, config, key_store)
		default:
			fmt.Printf("Unknown RESP indicator '%s'\n", string(tmp[0]))
		}
	}
}

func bulkArray(conn net.Conn, config Config, key_store KeyStore) {
	reader := bufio.NewReader(conn)

	number_of_elements_str, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading number of elements in bulk array")
		return
	}

	number_of_elements, err := strconv.Atoi(strings.TrimSpace(number_of_elements_str))
	if err != nil {
		fmt.Printf("Error parsing number of elements in bulk array: '%s'", number_of_elements_str)
	}
	elements := make([]interface{}, number_of_elements)

	fmt.Printf("Bulk array with %d elements\n", number_of_elements)

	for i := 0; i < number_of_elements; i++ {
		element, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		switch element[0] {
		case '$':
			element, err := bulkString(element[1:], reader)
			if err != nil {
				fmt.Printf(
					"Error parsing bulk string: %s",
					err)
			}
			elements[i] = element
		default:
			fmt.Printf("Unknown RESP type '%s'",
				element)
		}
	}

	command(conn, config, key_store, elements[0].(string), elements[1:])
}

func bulkString(length_str string, reader *bufio.Reader) (string, error) {
	string_length, err := strconv.Atoi(strings.TrimSpace(length_str))
	if err != nil {
		return "",
			fmt.Errorf(
				"Could not parse bulk string length '%s'", length_str)
	}

	str, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	str = strings.TrimSpace(str)

	if string_length != len(str) {
		return str, fmt.Errorf("Bulk string length %d mismatch with str val '%s' length %d",
			string_length, str, len(str))
	}

	return str, nil
}

func command(conn net.Conn, cfg Config, key_store KeyStore,
	command string, command_args []interface{}) {
	switch strings.ToLower(command) {
	case "set":
		set_args := make([]string, 0)

		for _, e := range command_args {
			set_args = append(set_args,
				e.(string))
		}
		set(conn, key_store, set_args)
	case "get":
		get_args := make([]string, 0)

		for _, e := range command_args {
			get_args = append(get_args,
				e.(string))
		}
		get(conn, key_store, get_args)
	case "echo":
		echo_args := make([]string, 0)

		for _, e := range command_args {
			echo_args = append(echo_args,
				e.(string))
		}
		echo(conn, echo_args)
	case "config":
		config_args := make([]string, 0)

		for _, e := range command_args {
			config_args = append(config_args,
				e.(string))
		}
		config(conn, cfg, config_args)

	case "ping":
		ping(conn)
	default:
		fmt.Printf("Unknown command %s", command)
	}
}
