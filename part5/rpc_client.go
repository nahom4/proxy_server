package main
import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"sync"
	"time"
)

type Backend struct {
    net.Conn
    Reader *bufio.Reader
    Writer *bufio.Writer
}

var backendQueue chan *Backend
var requestLock sync.Mutex

type Stats struct {
	RequestBytes map[string]int64
}

var stats Stats
var overallStats Stats
func init() {
	backendQueue = make(chan *Backend, 10)
	overallStats.RequestBytes = make(map[string]int64)
}

func getBackend() (*Backend, error) {
    select {
    case be := <-backendQueue:
        return be, nil
    case <-time.After(100 * time.Millisecond):
        be, err := net.Dial("tcp", "127.0.0.1:8081")
        if err != nil {
            return nil, err
        }
        return &Backend{
            Conn:   be,
            Reader: bufio.NewReader(be),
            Writer: bufio.NewWriter(be),
        }, nil
    }
}

func queueBackend(be *Backend) {
	select {
	case backendQueue <- be:
	case <-time.After(1 * time.Second):
		be.Close()
	}
}

func updateStats(req *http.Request, resp *http.Response) int64 {
	requestLock.Lock()
	defer requestLock.Unlock()

	client, err := rpc.DialHTTP("tcp", "localhost:8079")
	if err != nil {
		log.Fatalf("Failed to connect to rpc server: %s", err)
	}

	
	if err := client.Call("RpcServer.GetStats", &struct{}{}, &stats); err != nil {
		log.Fatalf("Failed to call GetStats: %s", err)
	}
	fmt.Println(stats.RequestBytes)
	fmt.Println(req.URL.Path,"URL PATH")
	path := req.URL.Path[1:]
	overallStats.RequestBytes[path] += stats.RequestBytes[path]
	fmt.Println(overallStats.RequestBytes[path])
	return overallStats.RequestBytes[path]
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		req, err := http.ReadRequest(reader)
		if err != nil {
			if err != io.EOF {
				log.Printf("Failed to read request: %s", err)
				return 
			}
		}
		
		if req == nil {
			log.Println("Received nil request")
			return
		}

		be, err := getBackend()
		if err != nil {
			log.Printf("Failed to get backend: %s", err)
			return
		}
	
		if err := req.Write(be.Writer); err != nil {
			log.Printf("Failed to write request: %s", err)
			return
		}
		
		be.Writer.Flush()
		if resp, err := http.ReadResponse(be.Reader, req); err == nil {
			bytes := updateStats(req, resp)
			resp.Header.Set("X-Bytes", strconv.FormatInt(bytes, 10))

			if err := resp.Write(conn); err == nil {
				log.Printf("%s : %d", req.URL.Path, resp.StatusCode)
			}
		}
		go queueBackend(be)
		
	}
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %s", err)
	}

	fmt.Println("Server is listening on port 8080...")
	for {
		if conn, err := ln.Accept(); err == nil {
			go handleConnection(conn)
		}
	}
}

