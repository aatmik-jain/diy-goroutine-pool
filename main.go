package main

import (
	"diy_goroutine_pool/pool"
	"time"
)

func main() {
	routinePool, err := pool.CreatePool(4, 12, 1001*time.Millisecond)
	if err != nil {
		return
	}
	routinePool.Run()
}
