// test_lock.go
package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/Vishal2002/kv_server/kvsrv"
	lock "github.com/Vishal2002/kv_server/locks"
)

func main() {
	clerk := kvsrv.MakeClerk("localhost:1234")

	// Test basic lock functionality
	fmt.Println("Testing lock functionality...")

	lock1 := lock.MakeLock(clerk, "test_lock")

	fmt.Println("Acquiring lock...")
	lock1.Acquire()
	fmt.Println("Lock acquired!")

	// Try to acquire with another lock instance (should wait)
	lock2 := lock.MakeLock(clerk, "test_lock")

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		fmt.Println("Second client trying to acquire...")
		lock2.Acquire()
		fmt.Println("Second client got lock!")
		time.Sleep(1 * time.Second)
		lock2.Release()
		fmt.Println("Second client released lock")
	}()

	time.Sleep(2 * time.Second)
	fmt.Println("First client releasing lock...")
	lock1.Release()
	fmt.Println("First client released lock")

	wg.Wait()
	fmt.Println("Test completed!")
}
