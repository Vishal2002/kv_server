package kvsrv

import "errors"

// Error types
var (
	ErrNoKey   = errors.New("key does not exist")
	ErrVersion = errors.New("version mismatch")
	ErrMaybe   = errors.New("operation may or may not have been executed")
)

// Error strings for RPC transmission
const (
	ErrNoKeyStr   = "key does not exist"
	ErrVersionStr = "version mismatch"
	ErrMaybeStr   = "operation may or may not have been executed"
)

// RPC request/response structures
type PutArgs struct {
	Key     string
	Value   string
	Version int
}

type PutReply struct {
	Err string // Changed from error to string
}

type GetArgs struct {
	Key string
}

type GetReply struct {
	Value   string
	Version int
	Err     string // Changed from error to string
}
