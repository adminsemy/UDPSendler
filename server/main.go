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

	// Read from UDP listener in endless loop
	countGorutines := runtime.NumCPU()
	readCh := make(chan DataFile)
	for i := 0; i < countGorutines; i++ {
		go readData(conn, readCh)
	}
	for {
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
		file, err := os.Create(string(fileName))
		if err != nil {
			slog.Error("Error create file:", "Err:", err)
			return
		}
		index := uint64(0)
		gotBite := int64(0)
		var bufferData []DataFile
		for data := range readCh {
			if data.Part < index {
				continue
			}
			if data.Part == index {
				file.Write(data.Body)
				index++
			} else {
				index = checkBuffer(bufferData, data, file, index)
				bufferData = append(bufferData, data)
			}
			gotBite += int64(len(data.Body))
			if gotBite == size {
				break
			}
		}
		fmt.Println("Done")
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
		slog.Info("New part:", "part:", string(data.Body))
		conn.WriteToUDP([]byte("OK\n"), addr)
		readCh <- data
	}
}

func checkBuffer(buf []DataFile, data DataFile, file *os.File, index uint64) uint64 {
	if len(buf) <= int(index) {
		return index
	}
	for i := 0; i < len(buf); i++ {
		if buf[i].Part == index {
			file.Write(buf[i].Body)
			buf[len(buf)-1], buf[i] = buf[i], buf[len(buf)-1]
			buf = buf[:len(buf)-1]
			index++
			i = -1
		}
	}
	return index
}

func main() {
	connection()
}
