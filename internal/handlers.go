package internal

import (
	"fmt"
	"net"
)

var DataBase map[string]string

func Ping(conn *net.UDPConn, addr *net.UDPAddr, packet string) {
	if len(packet) != 0 {
		_, err := conn.WriteToUDP([]byte("(error) ERR wrong number of arguments for 'PING' command\n"), addr)
		if err != nil {
			fmt.Println("Error sending response:", err)
		} else {
			fmt.Printf("Sent 'Error' to %s\n\n", addr)
		}
		return
	}

	_, err := conn.WriteToUDP([]byte("PONG\n"), addr)
	if err != nil {
		fmt.Println("Error sending response:", err)
	} else {
		fmt.Printf("Sent 'PONG' to %s\n\n", addr)
	}
}

func Set(conn *net.UDPConn, addr *net.UDPAddr, packet string) {
	if len(packet) < 2 {
		_, err := conn.WriteToUDP([]byte("(error) ERR wrong number of arguments for 'SET' command\n"), addr)
		if err != nil {
			fmt.Println("Error sending response:", err)
		} else {
			fmt.Printf("Sent 'Error' to %s\n\n", addr)
		}
		return
	}

	key, val := "", ""
	for i, v := range packet {
		if v == ' ' {
			val = packet[i+1:]
			break
		}

		key += string(v)
	}

	ok := false
	for _, v := range val {
		if v != ' ' {
			ok = true
		}
	}

	if !ok || len(val) == 0 {
		_, err := conn.WriteToUDP([]byte("(error) ERR wrong number of arguments for 'SET' command\n"), addr)
		if err != nil {
			fmt.Println("Error sending response:", err)
		} else {
			fmt.Printf("Sent 'Error' to %s\n\n", addr)
		}
		return
	}

	if DataBase == nil {
		DataBase = make(map[string]string) // Initialize the map if it's nil
	}
	DataBase[key] = val

	_, err := conn.WriteToUDP([]byte("OK\n"), addr)
	if err != nil {
		fmt.Println("Error sending response:", err)
	} else {
		fmt.Printf("Sent 'OK' to %s\n\n", addr)
	}
}

func Get(conn *net.UDPConn, addr *net.UDPAddr, packet string) {
	if DataBase == nil {
		DataBase = make(map[string]string) // Initialize the map if it's nil
	}
	val := DataBase[packet]
	if val == "" {
		_, err := conn.WriteToUDP([]byte("(nil)"), addr)
		if err != nil {
			fmt.Println("Error sending response:", err)
		} else {
			fmt.Printf("Sent '(nil)' to %s\n\n", addr)
		}
	}
	_, err := conn.WriteToUDP([]byte(val+"\n"), addr)
	if err != nil {
		fmt.Println("Error sending response:", err)
	} else {
		fmt.Printf("Sent 'value' to %s\n\n", addr)
	}
}
