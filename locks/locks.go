package lock

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Vishal2002/kv_server/client"
	"github.com/Vishal2002/kv_server/rpc"
)

type Lock struct {
	ck       *client.Clerk
	lockKey  string
	clientID string
}

func generateRandom() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(90000000) + 10000000
	return fmt.Sprintf("%d", code)
}

func MakeLock(ck *client.Clerk, lockKey string) *Lock {
	return &Lock{
		ck:       ck,
		lockKey:  lockKey,
		clientID: generateRandom(),
	}
}

func (lk *Lock) Acquire() {
	for {
		value, version, getErr := lk.ck.Get(lk.lockKey)

		if getErr == nil && value == lk.clientID {
			// We already hold the lock
			return
		}

		if getErr == rpc.ErrNoKey || (getErr == nil && value == "") {
			// Lock appears free; try to claim it
			putVersion := 0
			if getErr == nil {
				putVersion = version
			}
			putErr := lk.ck.Put(lk.lockKey, lk.clientID, putVersion)

			if putErr == nil {
				// Successfully acquired
				return
			}

			if putErr == rpc.ErrMaybe {
				// Ambiguous; check if we now hold it
				checkValue, _, checkErr := lk.ck.Get(lk.lockKey)
				if checkErr == nil && checkValue == lk.clientID {
					return
				}
			}
			// Else: race lost or error; retry loop will handle
		}

		// Lock held by someone else (or other error); wait and retry
		time.Sleep(50 * time.Millisecond)
	}
}

func (lk *Lock) Release() {
	for {
		// Get current lock state
		value, version, err := lk.ck.Get(lk.lockKey)

		if err == rpc.ErrNoKey {
			// Lock doesn't exist, nothing to release
			return
		}

		if err == nil {
			if value != lk.clientID {
				// Someone else has the lock, we can't release it
				return
			}

			// Try to release by putting empty value with current version
			putErr := lk.ck.Put(lk.lockKey, "", version)
			if putErr == nil || putErr == rpc.ErrMaybe {
				// Successfully released (or maybe released)
				return
			}
			// If we got ErrVersion, someone else might have changed it
			// Try again
		}

		time.Sleep(10 * time.Millisecond)
	}
}
