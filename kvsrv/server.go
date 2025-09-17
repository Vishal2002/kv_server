package kvsrv

import (
	"sync"
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

func (s *Server) Get(args *GetArgs, reply *GetReply) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if kv, ok := s.store[args.Key]; ok {
		reply.Value = kv.Value
		reply.Version = kv.Version
		reply.Err = "" // No error
	} else {
		reply.Err = ErrNoKeyStr
	}
	return nil
}

func (s *Server) Put(args *PutArgs, reply *PutReply) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.store[args.Key]

	if args.Version == 0 {
		// Creating new key
		if exists {
			reply.Err = ErrVersionStr
		} else {
			s.store[args.Key] = KeyValue{Value: args.Value, Version: 1}
			reply.Err = "" // No error
		}
		return nil
	}

	if !exists {
		reply.Err = ErrNoKeyStr
		return nil
	}

	if existing.Version != args.Version {
		reply.Err = ErrVersionStr
		return nil
	}

	s.store[args.Key] = KeyValue{Value: args.Value, Version: args.Version + 1}
	reply.Err = "" // No error
	return nil
}
