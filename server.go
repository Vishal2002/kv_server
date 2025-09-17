package server

import (
	"log"
	"net"
	"net/rpc"
	"sync"

	kvrpc "github.com/Vishal2002/kv_server/rpc"
)

type KeyValue struct {
	Value   string
	Version int
}

type Server struct {
	mu    sync.Mutex
	store map[string]KeyValue
}

func NewServer() *Server {
	return &Server{
		store: make(map[string]KeyValue),
	}
}

func (s *Server) Get(args *kvrpc.GetArgs, reply *kvrpc.GetReply) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if kv, ok := s.store[args.Key]; ok {
		reply.Value = kv.Value
		reply.Version = kv.Version
		reply.Err = nil
	} else {
		reply.Err = rpc.ErrNoKey
	}
	return nil
}

func (s *Server) Put(args *kvrpc.PutArgs, reply *kvrpc.PutReply) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.store[args.Key]

	if args.Version == 0 {
		// Creating new key
		if exists {
			reply.Err = rpc.ErrVersion
		} else {
			s.store[args.Key] = KeyValue{Value: args.Value, Version: 1}
			reply.Err = nil
		}
		return nil
	}

	if !exists {
		reply.Err = rpc.ErrNoKey
		return nil
	}

	if existing.Version != args.Version {
		reply.Err = rpc.ErrVersion
		return nil
	}

	s.store[args.Key] = KeyValue{Value: args.Value, Version: args.Version + 1}
	reply.Err = nil
	return nil
}

func main() {
	server := NewServer()
	rpc.RegisterName("KV_Server", server)

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
