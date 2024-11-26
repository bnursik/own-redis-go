package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"own-redis/internal"
	"own-redis/models"
	"strconv"
	"strings"
)

func main() {
	flag.Parse()

	if *models.HelpFlag {
		models.HelpMessage()
		os.Exit(0)
	}

	port, err := strconv.Atoi(*models.PortFlag)
	if err != nil || port < 1024 {
		log.Fatal("Incorrect port")
	}

	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:"+*models.PortFlag)
	if err != nil {
		fmt.Println("Error resolving address:", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Error listening on UDP address:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("UDP server is running on %s...\n", *models.PortFlag)

	buffer := make([]byte, 1024)

	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP connection:", err)
			continue
		}

		message := string(buffer[:n])
		if message == "\n" {
			continue
		}
		if string(message[len(message)-1]) == "\n" {
			message = message[:len(message)-1]
		}
		fmt.Printf("Received '%s' from %s\n", message, addr)

		command := ""
		packet := ""
		for i, v := range message {
			if v == ' ' {
				packet = string(message[i+1:])
				break
			}
			command += string(v)
		}

		switch strings.ToLower(command) {
		case "ping":
			internal.Ping(conn, addr, packet)
		case "set":
			internal.Set(conn, addr, packet)
		case "get":
			internal.Get(conn, addr, packet)
		default:
			_, err := conn.WriteToUDP([]byte("(error) ERR command not found!\n"), addr)
			if err != nil {
				fmt.Println("Error sending response:", err)
			} else {
				fmt.Printf("Sent 'Error' to %s\n\n", addr)
			}
		}
	}
}
