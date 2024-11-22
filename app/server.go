package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

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

		go handle(conn)

	}

}

func handle(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Handling connection")

	tmp := make([]byte, 1)

	_, err := conn.Read(tmp)
	if err != nil {
		fmt.Printf("Error reading first byte of request %s", err)
	}

	switch tmp[0] {
	case '*':
		bulkArray(conn)
	default:
		fmt.Printf("Unknown RESP indicator '%s'", string(tmp[0]))
	}
}

func bulkArray(conn net.Conn) {
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

	command(conn, elements[0].(string), elements[1:])
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

func command(conn net.Conn, command string, command_args []interface{}) {
	switch strings.ToLower(command) {
	case "echo":
		echo_args := make([]string, 0)

		for _, e := range command_args {
			echo_args = append(echo_args,
				e.(string))
		}
		echo(conn, echo_args)
	case "ping":
		ping(conn)
	default:
		fmt.Printf("Unknown command %s", command)
	}
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
	fmt.Println("writing PONG response")
	conn.Write([]byte("+PONG\r\n"))
}
