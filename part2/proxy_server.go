package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

func handleConnection(conn net.Conn) {
    defer conn.Close()
    reader := bufio.NewReader(conn)

    for {
        req, err := http.ReadRequest(reader)
        if err != nil {    
            if err != io.EOF {
                log.Printf("Failed to read request: %s", err)
            }
            return 
        }

        if be, err := net.Dial("tcp", "127.0.0.1:8081"); err == nil {
            be_reader := bufio.NewReader(be)
            if err := req.Write(be); err == nil {
                if resp, err := http.ReadResponse(be_reader, req); err == nil {
                    if err := resp.Write(conn); err == nil {
                        log.Printf("%s : %d", req.URL.Path, resp.StatusCode)
                    }
                }
            }
            be.Close() 
        }
    }
}
func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Failed to listen: %s",err)
	}

	fmt.Println("Server is listening on port 8080...")
	for {
		if conn, err := ln.Accept(); err == nil {

			go handleConnection(conn)
		}
	}
}

