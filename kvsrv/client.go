package kvsrv

import (
	"errors"
	"net/rpc"
	"time"
)

type Clerk struct {
	client *rpc.Client
}

func MakeClerk(server string) *Clerk {
	client, err := rpc.Dial("tcp", server)
	if err != nil {
		panic(err)
	}
	return &Clerk{client: client}
}

// Get method - simple, just call the server
func (ck *Clerk) Get(key string) (string, int, error) {
	args := &GetArgs{Key: key}
	reply := &GetReply{}

	// Keep retrying until we get a response
	for {
		err := ck.client.Call("KVServer.Get", args, reply)
		if err == nil {
			// Convert string error back to Go error
			if reply.Err == "" {
				return reply.Value, reply.Version, nil
			} else if reply.Err == ErrNoKeyStr {
				return reply.Value, reply.Version, ErrNoKey
			}
			// Handle other error strings if needed
			return reply.Value, reply.Version, errors.New(reply.Err)
		}
		// If RPC failed, wait a bit and retry
		time.Sleep(100 * time.Millisecond)
	}
}

// Put method - more complex due to at-most-once semantics
func (ck *Clerk) Put(key, value string, version int) error {
	args := &PutArgs{
		Key:     key,
		Value:   value,
		Version: version,
	}
	reply := &PutReply{}

	firstAttempt := true

	for {
		err := ck.client.Call("KVServer.Put", args, reply)
		if err == nil {

			if reply.Err == "" {
				return nil
			} else if reply.Err == ErrVersionStr && !firstAttempt {

				return ErrMaybe
			} else if reply.Err == ErrVersionStr {
				return ErrVersion
			} else if reply.Err == ErrNoKeyStr {
				return ErrNoKey
			} else if reply.Err == ErrMaybeStr {
				return ErrMaybe
			}
			// Handle unknown error strings
			return errors.New(reply.Err)
		}

		// RPC failed, need to retry
		firstAttempt = false
		time.Sleep(100 * time.Millisecond)
	}
}
