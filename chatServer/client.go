package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Connect error:", err)
		return
	}
	defer conn.Close()

	// Server se messages read
	go func() {
		reader := bufio.NewReader(conn)
		for {
			message, _ := reader.ReadString('\n')
			fmt.Print(message)
		}
	}()

	// User se input lekar server ko bhejo
	writer := bufio.NewWriter(conn)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		text := scanner.Text()
		fmt.Fprintln(writer, text)
		writer.Flush()
	}
}
