package server

import (
	"fmt"
	"net"
	"net-cat/pkg/utils"
	"strings"
	"sync"
)

type Server struct {
	address          string
	clients          map[net.Conn]string
	mu               sync.Mutex
	broadcast        chan string
	history          []string
	connectionsCount int
	connectionsMax   int
}

func NewServer(address string) *Server {
	return &Server{
		address:        address,
		clients:        make(map[net.Conn]string),
		broadcast:      make(chan string),
		connectionsMax: utils.ConnectionsLimit,
	}
}

func (s *Server) Start() error {
	address := s.address
	if !strings.Contains(address, ":") {
		address = ":" + address
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	addr := listener.Addr().String()
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		fmt.Println("Error extracting port:", err)
		return err
	}
	fmt.Printf("Listening on the port :%s\n", port)

	go s.handleMessages()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}

		s.mu.Lock()
		if s.connectionsCount >= s.connectionsMax {
			fmt.Println("Maximum connections reached, rejecting new connection")
			conn.Write([]byte("Server full, try again later.\n"))
			conn.Close()
			s.mu.Unlock()
			continue
		}
		s.connectionsCount++
		s.mu.Unlock()

		go s.handleConnection(conn)
	}
}

// Функция для отправки предыдущих сообщений новому клиенту
func sendPreviousMessages(conn net.Conn, history []string) error {
	for _, msg := range history {
		// Отправляем каждое сообщение из истории новому клиенту
		_, err := conn.Write([]byte(msg + "\n"))
		if err != nil {
			return fmt.Errorf("failed to send history message: %v", err)
		}
	}
	return nil
}
