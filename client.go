package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"net"
	"time"
)

// client connects to ClientServer, log-in, and put Conn into chan readyClientQueue
func newClient() {

	if config.ClientsReserve < 1 {
		config.ClientsReserve = 1
	}

	password := []byte(config.Password)
	attempts := 5
	idleClients := make(chan bool, config.ClientsReserve-1)
	defer close(idleClients)
	logger.Printf("Connecting to %s:%d\n", config.ServerHost, config.ServerClientPort)

	for {

		if attempts < 0 {
			logger.Panic("too many failed attempts. Something is wrong with server-client dialog")
		}

		// make Conn for connections to the server
		client, err := net.Dial("tcp", fmt.Sprintf("%s:%d", config.ServerHost, config.ServerClientPort))
		client.SetDeadline(time.Now().Add(time.Millisecond * 1000))
		if err != nil {
			logger.Panic("could not connect to ClientServer")
		}

		// Authorisation
		challenge, err := readConn(client, challengeSize)
		if err != nil {
			logger.Println("Server disconected on receiving challenge")
			client.Close()
			attempts -= 1
			continue
		}
		controllPhrase := md5.Sum(bytes.Join([][]byte{challenge, password}, []byte("")))
		_, err = client.Write(controllPhrase[:])
		if err != nil {
			logger.Println("Server disconected on sending controll phrase")
			client.Close()
			attempts -= 1
			continue
		}

		status, err := readConn(client, len(statusAccepted))
		if err != nil {
			logger.Println("Server disconected on receiving status")
			client.Close()
			attempts -= 1
			continue
		}

		client.SetDeadline(time.Time{})
		if bytes.Equal(status[:], []byte(statusAccepted)) { // Client accepted
			go clientActivate(client, idleClients)
		}

		// blocks while there are less than config.ClientsReserve idle clients
		// clients read one when they activating
		idleClients <- true
		attempts = 5
	}
}

// runClient waits for startCommand from client.
// makes Conn to servise and launches transfer loop
func clientActivate(client net.Conn, report <-chan bool) {
	// Waiting for START command
	defer client.Close()
	response, err := readConn(client, len(startCommand))
	<-report
	if err != nil {
		return
	}

	// Check if START command is correct
	if !bytes.Equal([]byte(startCommand), response[:]) {
		logger.Printf("Received wrong start command: %s\n", string(response))
		return
	}

	// Connect to service.
	service, err := net.Dial("tcp", fmt.Sprintf("%s:%d", config.ServiceHost, config.ServicePort))
	if err != nil {
		logger.Panic("could not connect to the Service")
	}
	defer service.Close()

	go pump(client, service)
	pump(service, client)

}
