package main

import (
	"diy_goroutine_pool/pool"
	"testing"
	"time"
)

func TestSingleton(t *testing.T) {
	pool1, _ := pool.CreatePool(1, 3, 2*time.Second)
	pool2, _ := pool.CreatePool(4, 11, 1000*time.Millisecond)
	//now pool1 should have numRoutines=4 and timeout=1001*time.Millisecond too

	if pool1 != pool2 {
		t.Errorf("Pool instances are different")
	}
	pool2.Run()
}

func TestExcessJobs(t *testing.T) {
	_, err := pool.CreatePool(5, 20, 2*time.Second)

	if err == nil {
		t.Errorf("Should not have create more than 15 jobs")
	}
}
