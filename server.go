package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"net"
	"time"
)

// newServer launches serverClientFactory for clientQueue channel
// Then newServer combine each public connection with a client and starts new pump for them
func newServer() {
	clientQueue := make(chan net.Conn, 100)
	defer close(clientQueue)

	// launching client maker
	go serverClientFactory(clientQueue)

	// make Conn for public to connect to
	logger.Printf("Launching publicServer on port %d\n", config.ServerPublicPort)
	publicServer, err := net.Listen("tcp", fmt.Sprintf(":%d", config.ServerPublicPort))
	if err != nil {
		logger.Fatalf("could not listen on %d port\n", config.ServerPublicPort)
	}
	defer publicServer.Close()
	for {
		// Pulling new client
		client := <-clientQueue
		// Accepting new Public
		public, err := publicServer.Accept()
		if err != nil {
			public.Close()
			continue
		}
		// Start pumping
		client.Write([]byte(startCommand))
		go pump(public, client)
		go pump(client, public)
	}
}

// serverClientFactory listens for clients, authentificate and sends them into clientQueue channel.
// Method of authentification - Digest
func serverClientFactory(clientQueue chan<- net.Conn) {
	password := []byte(config.Password)
	// make Conn for clients to connect to
	logger.Printf("Launching clientServer on port %d\n", config.ServerClientPort)
	clientServer, err := net.Listen("tcp", fmt.Sprintf(":%d", config.ServerClientPort))
	if err != nil {
		logger.Fatalf("could not listen on %d port\n", config.ServerClientPort)
	}
	defer clientServer.Close()

	for {
		// Accepting new Client
		client, err := clientServer.Accept()
		if err != nil {
			client.Close()
			continue
		}

		clientHost, _, _ := net.SplitHostPort(client.RemoteAddr().String())
		if isHostBanned(clientHost) {
			client.Close()
			continue
		}
		challenge := GetRandStringBytes(challengeSize)
		client.SetDeadline(time.Now().Add(time.Millisecond * 1000))

		// sending secret
		_, err = client.Write(challenge)
		if err != nil {
			client.Close()
			banHost(clientHost)
			logger.Printf("BAN:%s - disconected on log-in\n", clientHost)
			continue
		}

		// reading response
		response, err := readConn(client, 16)
		if err != nil {
			client.Close()
			banHost(clientHost)
			logger.Printf("BAN:%s - disconected on log-in\n", clientHost)
			continue
		}

		// check if password is correct
		controllPhrase := md5.Sum(bytes.Join([][]byte{challenge, password}, []byte("")))
		if !bytes.Equal(controllPhrase[:], response[:]) { // wrong password
			client.Close()
			banHost(clientHost)
			logger.Printf("BAN:%s - wrong password", clientHost)
			continue
		}

		// Client is authorized
		// Send OK status
		// Put it into channel
		_, _ = client.Write([]byte(statusAccepted))
		client.SetDeadline(time.Time{})
		clientQueue <- client
		unbanExpired()
	}
}
