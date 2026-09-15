package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

type Client struct {
	conn net.Conn
	name string
}

type ChatServer struct {
	clients    map[*Client]bool
	broadcast  chan string
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan string),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (s *ChatServer) Run() {
	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			s.clients[client] = true
			s.mu.Unlock()
			s.broadcast <- fmt.Sprintf("%s joined the chat!", client.name)

		case client := <-s.unregister:
			s.mu.Lock()
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
			}
			s.mu.Unlock()
			s.broadcast <- fmt.Sprintf("%s left the chat.", client.name)

		case message := <-s.broadcast:
			s.mu.Lock()
			for client := range s.clients {
				fmt.Fprintf(client.conn, "%s\n", message)
			}
			s.mu.Unlock()
		}
	}
}

func handleClient(conn net.Conn, server *ChatServer) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Client se naam poocho
	fmt.Fprintf(conn, "Enter your name: ")
	name, _ := reader.ReadString('\n')
	name = name[:len(name)-1] // newline remove

	client := &Client{conn: conn, name: name}
	server.register <- client

	// Messages read karo
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			server.unregister <- client
			break
		}

		formatted := fmt.Sprintf("[%s]: %s", client.name, message[:len(message)-1])
		server.broadcast <- formatted
	}
}

func main() {
	server := NewChatServer()
	go server.Run()

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Listen error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Chat server started on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}

		go handleClient(conn, server)
	}
}
