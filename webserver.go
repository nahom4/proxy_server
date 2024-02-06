package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)


var requestLock sync.Mutex

type Stats struct {
	RequestBytes map[string]int64
}

type Empty struct{}
type RpcServer struct{}

var stats Stats

func init() {
    stats.RequestBytes = make(map[string]int64)
    stats.RequestBytes["localhost:8037"] = 10
}


func (r *RpcServer) GetStats(args *Empty, reply *Stats) error {
	requestLock.Lock()
	defer requestLock.Unlock()
	reply.RequestBytes = make(map[string]int64)
	for k, v := range stats.RequestBytes {
		reply.RequestBytes[k] = v
	}

	return nil
}

func documentationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	http.Redirect(w, r, "https://pkg.go.dev/", http.StatusSeeOther)
}

func main() {
	rpc.Register(&RpcServer{})
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":8079")
	if err != nil {
		log.Fatalf("Failed to listen: %s", err)
	}
	go http.Serve(l, nil)
	http.HandleFunc("/", documentationHandler)
	log.Println("Server is listening on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
