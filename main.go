package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
)

var config = new(struct {
	Password         string `json:"password"`
	Mode             string `json:"mode"`        // service|client
	ServerHost       string `json:"serverHost"`  // server address
	ServiceHost      string `json:"serviceHost"` // service address
	ServerPublicPort int    `json:"serverPublicPort"`
	ServerClientPort int    `json:"serverClientPort"`
	ServicePort      int    `json:"servicePort"`
	ClientsReserve   int    `json:"clientsReserve"`
})

var logger = log.New(os.Stdout, "", log.LstdFlags)

const (
	startCommand   string = "START"
	statusAccepted string = "OK"
	challengeSize  int    = 64
)

func main() {
	configFile := "./config.json"
	usage := fmt.Sprintf("Usage: %s PATH_TO_CONFIG_FILE\n default is %s ./config.json", os.Args[0], os.Args[0])
	if len(os.Args) == 2 {
		configFile = os.Args[1]
	}
	// read config.json
	s, err := os.ReadFile(configFile)
	if err != nil {
		err = fmt.Errorf("read file.json: %s", err)
		logger.Panic(err)
	}
	err = json.Unmarshal(s, &config)
	if err != nil {
		err = fmt.Errorf("unmarshal json: %s", err)
		logger.Panic(err)
	}

	switch config.Mode {
	case "server":
		newServer()
	case "client":
		newClient()
	default:
		logger.Println(usage)
		return
	}
}

// GetRandStringBytes returns random []bytes string of n bytes
func GetRandStringBytes(n int) []byte {
	symbols := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890.!#")
	b := make([]byte, n)
	for i := range b {
		b[i] = symbols[rand.Int63()%int64(len(symbols))]
	}
	return b
}

// transfer reads from one net.Conn and writes the data to another net.Conn
// returns when any error occured
// closes connections on returning
func pump(fromConn net.Conn, toConn net.Conn) {
	defer fromConn.Close()
	defer toConn.Close()
	buf := make([]byte, 2048)
	for {
		n, err := fromConn.Read(buf[:])
		if n == 0 && err == io.EOF {
			break
		}
		_, err = toConn.Write(buf[0:n])
		if err != nil {
			break
		}
	}
}

// readConn reads size bytes from conn and return them
func readConn(conn net.Conn, size int) (buffer []byte, err error) {
	buffer = make([]byte, size)
	for sz := 0; sz < size; {
		n, err := conn.Read(buffer[sz:])
		sz += n
		if n == 0 && err == io.EOF {
			conn.Close()
			return buffer[:], fmt.Errorf("disconnected on read")
		}
	}
	return buffer[:], nil
}
