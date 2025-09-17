package client

import (
	"net/rpc"
	"time"

	kvrpc "github.com/Vishal2002/kv_server/rpc"
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
	args := &kvrpc.GetArgs{Key: key}
	reply := &kvrpc.GetReply{}

	// Keep retrying until we get a response
	for {
		ok := ck.client.Call("KV_Server.Get", args, reply)
		if ok {
			return reply.Value, reply.Version, reply.Err
		}
		// If RPC failed, wait a bit and retry
		time.Sleep(100 * time.Millisecond)
	}
}

// Put method - more complex due to at-most-once semantics
func (ck *Clerk) Put(key, value string, version int) error {
	args := &kvrpc.PutArgs{
		Key:     key,
		Value:   value,
		Version: version,
	}
	reply := &kvrpc.PutReply{}

	firstAttempt := true

	for {
		ok := ck.client.Call("KV_Server.Put", args, reply)
		if ok {
			// Got a response
			if reply.Err == rpc.ErrVersion && !firstAttempt {
				// This is the tricky case - we retried and got ErrVersion
				// We can't tell if our first attempt succeeded or not
				return rpc.ErrMaybe
			}
			return reply.Err
		}

		// RPC failed, need to retry
		firstAttempt = false
		time.Sleep(100 * time.Millisecond)
	}
}
