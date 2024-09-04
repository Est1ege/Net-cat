package server

import (
	"bufio"
	"fmt"
	"net"
	"net-cat/pkg/utils"
	"strings"
	"time"
)

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Client connected:", conn.RemoteAddr())

	_, err := conn.Write([]byte(utils.WelcomeMsg))
	if err != nil {
		fmt.Println("Failed to send welcome message to client:", err)
		return
	}

	buffer, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read from client:", err)
		return
	}
	userName := strings.TrimSpace(buffer)
	validName, invalidNameErrCode := utils.IsValidName(userName, s.clients)
	if !validName {
		errMsg := "Invalid username. "
		if invalidNameErrCode == 1 {
			errMsg = errMsg + "Usernames cannot contain spaces and must be longer than 1 symbol.\n"
		} else if invalidNameErrCode == 2 {
			errMsg = errMsg + "This username has been taken, try different one.\n"
		}
		conn.Write([]byte(errMsg))
		return
	}

	s.mu.Lock()
	s.clients[conn] = userName
	historyCopy := make([]string, len(s.history))
	copy(historyCopy, s.history)
	s.mu.Unlock()

	// Отправка истории сообщений новому клиенту
	if err := sendPreviousMessages(conn, historyCopy); err != nil {
		fmt.Println("Failed to send previous messages to client:", err)
		return
	}

	s.broadcast <- fmt.Sprintf("%s has joined our chat...", userName)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := strings.TrimSpace(scanner.Text())
		if message == "" {
			continue
		}
		fmt.Printf("[%s] %s: %s\n", time.Now().Format("2006-01-02 15:04:05"), userName, message)
		s.broadcast <- fmt.Sprintf("[%s] %s: %s", time.Now().Format("2006-01-02 15:04:05"), userName, message)
	}

	s.mu.Lock()
	delete(s.clients, conn)
	s.connectionsCount--
	s.mu.Unlock()
	s.broadcast <- fmt.Sprintf("%s has left our chat...", userName)
}

func (s *Server) handleMessages() {
	for msg := range s.broadcast {
		s.mu.Lock()
		s.history = append(s.history, msg)

		for client := range s.clients {
			_, err := client.Write([]byte(msg + "\n"))
			if err != nil {
				fmt.Printf("Failed to send message to %s, removing from clients list\n", s.clients[client])
				client.Close()
				delete(s.clients, client)
			}
		}
		s.mu.Unlock()
	}
}
