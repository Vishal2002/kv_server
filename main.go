package main

import (
	"log"
	"net"
	"net/rpc"

	"github.com/Vishal2002/kv_server/kvsrv"
)

func main() {
	server := kvsrv.NewServer()
	rpc.RegisterName("KVServer", server)

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("Listen error:", err)
	}

	log.Println("KV Server listening on :1234")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}
		go rpc.ServeConn(conn)
	}
}
