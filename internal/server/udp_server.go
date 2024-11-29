package server

import (
	"log"
	"net"
	"own-redis/internal/storage"
	"strconv"
	"strings"
)

type UDPServer struct {
	port  int
	store *storage.Store
}

func NewUDPServer(port int, store *storage.Store) *UDPServer {
	return &UDPServer{
		port:  port,
		store: store,
	}
}

func (s *UDPServer) ListenAndServe() error {
	addr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(s.port))
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Printf("Listening on UDP port %d", s.port)

	for {
		buffer := make([]byte, 1024)
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading UDP packet: %v", err)
			continue
		}

		if string(buffer[:n]) == "\n" {
			continue
		}

		go s.handleRequest(conn, remoteAddr, buffer[:n])
	}
}

func (s *UDPServer) handleRequest(conn *net.UDPConn, addr *net.UDPAddr, data []byte) {
	request := strings.TrimSpace(string(data))
	parts := strings.Fields(strings.ToUpper(request))

	var response string
	switch {
	case len(parts) == 0:
		response = "(error) ERR invalid request"

	case parts[0] == "PING":
		response = s.store.Ping()

	case parts[0] == "SET":
		response = s.handleSetCommand(parts)

	case parts[0] == "GET":
		response = s.handleGetCommand(parts)

	default:
		response = "(error) ERR unknown command"
	}

	conn.WriteToUDP([]byte(response+"\n"), addr)
}

func (s *UDPServer) handleSetCommand(parts []string) string {
	if len(parts) < 3 {
		return "(error) ERR wrong number of arguments for 'SET' command"
	}

	havePx := false
	pxInd := 0
	for i, v := range parts {
		if v == "PX" {
			havePx = true
			pxInd = i
		}
	}

	var key string
	var value []string
	expireMs := int64(0)

	if havePx {
		if len(parts) < 5 {
			return "(error) ERR invalid expire time"
		}
		key, value = parts[1], parts[2:pxInd]

		if len(parts[pxInd+1:]) != 1 {
			return "(error) ERR invalid expire time"
		}

		ms, err := strconv.ParseInt(parts[pxInd+1], 10, 64)
		if err != nil || ms <= 0 {
			return "(error) ERR invalid expire time"
		}
		expireMs = ms
	} else {
		key = parts[1]
		value = parts[2:]
	}

	val := ""
	for i, v := range value {
		if i == len(value)-1 {
			val += v
		} else {
			val += (v + " ")
		}
	}

	s.store.Set(key, val, expireMs)
	return "OK"
}

func (s *UDPServer) handleGetCommand(parts []string) string {
	if len(parts) != 2 {
		return "(error) ERR wrong number of arguments for 'GET' command"
	}

	value, exists := s.store.Get(parts[1])
	if !exists {
		return "(nil)"
	}
	return value
}
