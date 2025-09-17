package rpc

import "errors"

// Error types
var (
	ErrNoKey   = errors.New("key does not exist")
	ErrVersion = errors.New("version mismatch")
	ErrMaybe   = errors.New("operation may or may not have been executed")
)

// RPC request/response structures
type PutArgs struct {
	Key     string
	Value   string
	Version int
}

type PutReply struct {
	Err error
}

type Server struct {
}

type GetArgs struct {
	Key string
}

type GetReply struct {
	Value   string
	Version int
	Err     error
}
