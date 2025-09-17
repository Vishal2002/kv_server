package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Vishal2002/kv_server/kvsrv"
)

func main() {
	fmt.Println("=== KV Server Throughput Benchmark ===")

	clients := []int{1, 5, 10, 20}
	duration := 5 * time.Second

	for _, numClients := range clients {
		fmt.Printf("\n--- Testing with %d clients ---\n", numClients)

		var operations int64
		var wg sync.WaitGroup
		start := time.Now()

		for i := 0; i < numClients; i++ {
			wg.Add(1)
			go func(clientID int) {
				defer wg.Done()
				clerk := kvsrv.MakeClerk("localhost:1234")

				for time.Since(start) < duration {
					key := fmt.Sprintf("bench-%d-%d", clientID, atomic.LoadInt64(&operations))

					// Put operation
					clerk.Put(key, "test-value", 0)
					atomic.AddInt64(&operations, 1)

					// Get operation
					clerk.Get(key)
					atomic.AddInt64(&operations, 1)
				}
			}(i)
		}

		wg.Wait()
		elapsed := time.Since(start)
		totalOps := atomic.LoadInt64(&operations)

		fmt.Printf("Operations: %d\n", totalOps)
		fmt.Printf("Time: %v\n", elapsed)
		fmt.Printf("Throughput: %.2f ops/sec\n", float64(totalOps)/elapsed.Seconds())
	}
}
