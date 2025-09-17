package main

import (
	"fmt"
	"time"

	"github.com/Vishal2002/kv_server/kvsrv"
)

func main() {
	// Wait a moment for server to be ready
	time.Sleep(1 * time.Second)

	// Create a clerk (client)
	clerk := kvsrv.MakeClerk("localhost:1234")

	// Test basic operations
	fmt.Println("Testing KV operations...")

	// Test Put
	err := clerk.Put("hello", "world", 0)
	fmt.Printf("Put 'hello'='world': %v\n", err)

	// Test Get
	value, version, err := clerk.Get("hello")
	fmt.Printf("Get 'hello': value=%s, version=%d, err=%v\n", value, version, err)

	// Test updating
	if err == nil {
		err = clerk.Put("hello", "updated", version)
		fmt.Printf("Update 'hello'='updated': %v\n", err)

		// Get again
		value, version, err = clerk.Get("hello")
		fmt.Printf("Get 'hello' again: value=%s, version=%d, err=%v\n", value, version, err)
	}
}
