package utils

import (
	"net"
	"strings"
)

func IsValidName(name string, clients map[net.Conn]string) (bool, int) {

	if len(name) < 1 {
		return false, 1
	}
	// Check if the name contains any space
	if strings.Contains(name, " ") {
		return false, 1
	}

	for _, clientName := range clients {
		if name == clientName {
			return false, 2 // Name is already taken
		}
	}

	return true, 0
}

func IsValidPort(port string) bool {
	//Проверка что строка не пуста
	if strings.TrimSpace(port) == "" {
		return false
	}

	// Преобразуйте порт в целое число
	portNum, err := parsePort(port)
	if err != true {
		return false
	}

	// Проверьте, что порт находится в допустимом диапазоне
	if portNum < 1 || portNum > 65535 {
		return false
	}

	return true
}

func parsePort(portStr string) (int, bool) {
	port := 0
	for _, char := range portStr {
		if char < '0' || char > '9' {
			return 0, false // Contains a non-digit character
		}
		port = port*10 + int(char-'0')
	}
	return port, true
}
