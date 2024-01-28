package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Resolve the string address to a UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Dial to the address with UDP
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Open file
	file, err := os.Open("test.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	info, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	l := info.Size()
	b := []byte{
		byte(0xff & l),
		byte(0xff & (l >> 8)),
		byte(0xff & (l >> 16)),
		byte(0xff & (l >> 24)),
		byte(0xff & (l >> 32)),
		byte(0xff & (l >> 40)),
		byte(0xff & (l >> 48)),
		byte(0xff & (l >> 56)),
	}
	// Send a message to the server
	_, err = conn.Write(append(b, []byte(file.Name())...))
	fmt.Println("send...")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	data, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))
	if string(data) == "OK" {
		fmt.Println("OK")
	}
}
