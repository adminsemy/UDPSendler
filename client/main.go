package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
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
	defer file.Close()
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

	if string(data) == "OK\n" {
		wg := sync.WaitGroup{}
		buf := make([]byte, 499)
		part := int64(0)
		var n int
		for {
			for i := 0; i < 4; i++ {
				n, err = file.Read(buf)
				if err == io.EOF {
					break
				}
				wg.Add(1)
				send := make([]byte, n)
				copy(send, buf)
				go func(part int64) {
					sendBytes(send, part, conn)
					wg.Done()
				}(part)
				part++
			}
			wg.Wait()
			fmt.Println("done")
			if err == io.EOF {
				break
			}
		}
	}
}

func sendBytes(send []byte, part int64, conn *net.UDPConn) {
	partBinary := make([]byte, 8)
	binary.LittleEndian.PutUint64(partBinary, uint64(part))
	data := append(partBinary, send...)
	for {
		_, err := conn.Write(data)
		if err != nil {
			fmt.Println(err)
			return
		}
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		if string(data) == "OK\n" {
			return
		}
	}
}
