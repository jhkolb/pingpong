package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

const numIters = 100
const delay = 1 * time.Second

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <ping | pong>\n", os.Args[0])
		os.Exit(1)
	}
	mode := os.Args[1]
	if mode != "ping" && mode != "pong" {
		fmt.Printf("Usage: %s <ping | pong>\n", os.Args[0])
		os.Exit(1)
	}

	if mode == "ping" {
		clientConn, err := net.Dial("tcp", "spawnpoint_pong:4590")
		if err != nil {
			fmt.Println("Failed to connect:", err)
			os.Exit(1)
		}

		sendMsg := "ping\n"
		reader := bufio.NewReader(clientConn)
		for i := 0; i < numIters; i++ {
			fmt.Printf("(%d) Sending: %s", i, sendMsg)
			if _, err := clientConn.Write([]byte(sendMsg)); err != nil {
				fmt.Println("Send failed:", err)
				os.Exit(1)
			}

			recvMsg, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Receive failed:", err)
				os.Exit(1)
			}
			fmt.Printf("(%d) Received: %s", i, recvMsg)

			time.Sleep(delay)
		}

		clientConn.Close()
	} else {
		ln, err := net.Listen("tcp", ":4590")
		if err != nil {
			fmt.Println("Listen failed:", err)
			os.Exit(1)
		}

		serverConn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept failed:", err)
			os.Exit(1)
		}
		reader := bufio.NewReader(serverConn)

		sendMsg := "pong\n"
		for i := 0; i < numIters; i++ {
			recvMsg, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Receive failed:", err)
				os.Exit(1)
			}
			fmt.Printf("(%d) Received: %s", i, recvMsg)

			fmt.Printf("(%d) Sending: %s\n", i, sendMsg)
			if _, err := serverConn.Write([]byte(sendMsg)); err != nil {
				fmt.Println("Send failed:", err)
				os.Exit(1)
			}
		}

		serverConn.Close()
	}
}
