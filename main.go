package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:1080")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		client, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleClient(client)
	}
}

func handleClient(client net.Conn) {
	defer client.Close()

	if err := authenticate(client); err != nil {
		log.Println("Authentication error:", err)
		return
	}

	target, err := getTargetAddress(client)
	if err != nil {
		log.Println("Error getting target address:", err)
		return
	}

	destConn, err := net.Dial("tcp", target)
	if err != nil {
		log.Println("Error connecting to target:", err)
		sendReply(client, 0x04) // Host unreachable
		return
	}
	defer destConn.Close()

	sendReply(client, 0x00) // Success

	go io.Copy(destConn, client)
	io.Copy(client, destConn)
}

func authenticate(client net.Conn) error {
	version := make([]byte, 1)
	if _, err := io.ReadFull(client, version); err != nil {
		return err
	}

	if version[0] != 5 {
		return fmt.Errorf("unsupported SOCKS version: %d", version[0])
	}

	nmethods := make([]byte, 1)
	if _, err := io.ReadFull(client, nmethods); err != nil {
		return err
	}

	methods := make([]byte, nmethods[0])
	if _, err := io.ReadFull(client, methods); err != nil {
		return err
	}

	// We only support no authentication (0x00)
	if _, err := client.Write([]byte{0x05, 0x00}); err != nil {
		return err
	}

	return nil
}

func getTargetAddress(client net.Conn) (string, error) {
	header := make([]byte, 4)
	if _, err := io.ReadFull(client, header); err != nil {
		return "", err
	}

	if header[0] != 0x05 || header[1] != 0x01 || header[2] != 0x00 {
		return "", fmt.Errorf("invalid SOCKS5 request")
	}

	var target string

	switch header[3] {
	case 0x01: // IPv4
		addr := make([]byte, 4)
		if _, err := io.ReadFull(client, addr); err != nil {
			return "", err
		}
		target = net.IP(addr).String()
	case 0x03: // Domain name
		domainLen := make([]byte, 1)
		if _, err := io.ReadFull(client, domainLen); err != nil {
			return "", err
		}
		domain := make([]byte, domainLen[0])
		if _, err := io.ReadFull(client, domain); err != nil {
			return "", err
		}
		target = string(domain)
	case 0x04: // IPv6
		addr := make([]byte, 16)
		if _, err := io.ReadFull(client, addr); err != nil {
			return "", err
		}
		target = net.IP(addr).String()
	default:
		return "", fmt.Errorf("unsupported address type")
	}

	port := make([]byte, 2)
	if _, err := io.ReadFull(client, port); err != nil {
		return "", err
	}
	portNum := binary.BigEndian.Uint16(port)

	return fmt.Sprintf("%s:%d", target, portNum), nil
}

func sendReply(client net.Conn, status byte) {
	reply := []byte{0x05, status, 0x00, 0x01, 0, 0, 0, 0, 0, 0}
	client.Write(reply)
}
