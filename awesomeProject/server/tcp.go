package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func handleConnection (conn net.Conn) {
	defer conn.Close()

	//	conn.SetDeadline(time.Now().Add(5* time.Second))
	log.Println("Reading the request")
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}  else {
			fmt.Println(line)
		}
	}

	body := `<!DOCTYPE html><html lang="en"><head><meta charet="UTF-8"><title></title></head><body>Hello world</body></html>`
	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	fmt.Fprint(conn, "\r\n")
	fmt.Fprint(conn, body)
	conn.Write([]byte("\r\n"))
	log.Println( "Server responded")
}

func handleSocketConnection (conn net.Conn) {
	defer conn.Close()

	log.Println("SocketServer: Reading the request")
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}  else {
			log.Println("SocketServer:" + line + "\r\n")
			_, _ = conn.Write([]byte("FromSocket: " + (line) + "\r\n"))
		}
	}
}

func StartEchoSocket() bool {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic("Cannot listen on 8081")
	}
	log.Println("Echo Socket server started")

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic("Cannot accept on 8081")
		}
		go handleSocketConnection(conn)
	}
	log.Println("Closing Echo Socket server")
	return true
}

func StartHTTPServer() bool {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic("Cannot listen on 8080")
	}
	log.Println("HTTP Server started")

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic("Cannot accept on 8080")
		}
		go handleConnection(conn)
	}
	log.Println("Closing HTTP server")
	return true
}