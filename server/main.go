package main

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"net"
	"os"
	"runtime"
)

type DataFile struct {
	Part uint64 `json:"part"`
	Body []byte `json:"body"`
}

func connection() {
	// Resolve the string address to a UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Start listening for UDP packages on the given address
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var buf [264]byte
	n, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		fmt.Println(err)
		return
	}
	size := int64(binary.LittleEndian.Uint64(buf[0:8]))
	fileName := buf[8:n]
	slog.Info("New file:", "name:", string(fileName), "size:", size)
	conn.WriteToUDP([]byte("OK\n"), addr)

	// Read from UDP listener in endless loop
	countGorutines := runtime.NumCPU()
	readCh := make(chan DataFile)
	for i := 0; i < countGorutines; i++ {
		go readData(conn, readCh)
	}
	index := uint64(0)
	for {
		var buf []DataFile
		var file *os.File
		for data := range readCh {
			if data.Part < index {
				continue
			}
			if data.Part == index {
				if index == 0 {
					file, _ = os.Create(string(fileName))
				}
				file.Write(data.Body)
				index++
			}
			buf = append(buf, data)
		}
		file.Close()
	}
}

func readData(conn *net.UDPConn, readCh chan<- DataFile) {
	var buf [507]byte
	var data DataFile
	for {
		n, addr, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			slog.Error("Error from read data:", "Err:", err)
			return
		}
		data.Part = binary.LittleEndian.Uint64(buf[0:8])
		data.Body = buf[8:n]
		readCh <- data
		conn.WriteToUDP([]byte("OK\n"), addr)
	}
}

func main() {
	connection()
}
