package pool

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"time"
)

type Pool interface {
	Run()
}

type pool struct {
	numRoutines int
	numJobs     int
	timeout     time.Duration
}

type result struct {
	jobID     int
	isDropped bool
}

var poolInstance *pool

func (p *pool) worker(id int, jobs <-chan int, results chan<- result) {
	for job := range jobs {

		ctx, _ := context.WithTimeout(context.Background(), p.timeout)
		temp := make(chan string)
		go func() {
			fmt.Println("worker", id, "started job:", job)
			time.Sleep(time.Second)
			temp <- "worker " + strconv.Itoa(id) + " finished job: " + strconv.Itoa(job)
		}()

		select {
		case <-ctx.Done():
			fmt.Println("Dropping job:", job, "due to timeout")
			results <- result{jobID: job, isDropped: true}
		case t := <-temp:
			fmt.Println(t)
			results <- result{jobID: job, isDropped: false}
		}
	}
}

func CreatePool(numRoutines int, numJobs int, timeout time.Duration) (Pool, error) {

	numCPUs := runtime.NumCPU()
	if numJobs > numCPUs/4*numRoutines {
		return nil, fmt.Errorf("cannot have more than %d jobs for %d routines", numCPUs/4*numRoutines, numRoutines)
	}

	if poolInstance == nil {
		poolInstance = new(pool)
	}
	poolInstance.numRoutines = numRoutines
	poolInstance.numJobs = numJobs
	poolInstance.timeout = timeout
	return poolInstance, nil
}

func (p *pool) Run() {
	jobs := make(chan int, p.numRoutines)
	results := make(chan result, p.numRoutines)

	for id := 1; id <= p.numRoutines; id++ {
		go p.worker(id, jobs, results)
	}

	for job := 1; job <= poolInstance.numJobs; job++ {
		jobs <- job
	}
	close(jobs)

	for i := 0; i < poolInstance.numJobs; i++ {
		res := <-results
		if res.isDropped {
			fmt.Println("Result for job", res.jobID, ": dropped")
		} else {
			fmt.Println("Result for job", res.jobID, ": completed")
		}
	}
}
