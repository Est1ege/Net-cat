package utils

import (
	"fmt"
	"log"
	"net-cat/internal/server"
	"net-cat/pkg/utils"
	"os"
)

func Run() {
	port := utils.DefaultPort // Порт по умолчанию

	if len(os.Args) == 2 {
		port = os.Args[1]
	} else if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	if !utils.IsValidPort(port) {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	srv := server.NewServer(port) // Запуск сервера
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
