package main

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Job struct {
	id           int
	randomNumber int
}

type Result struct {
	job         Job
	sumOfDigits int
}

// var (
// 	jobs    = make(chan Job, 10)
// 	results = make(chan Result, 10)
// )

func digits(number int) int {
	sum := 0
	no := number
	for no != 0 {
		digit := no % 10
		sum += digit
		no /= 10
	}
	return sum
}

func workerReal(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, jobs <-chan Job, results chan<- Result) {
	defer func() {
		wg.Done()
		log.Println("cancelled")
	}()

	for job := range jobs {
		select {
		case <-ctx.Done():
			return
		default:
		}
		output := Result{job, digits(job.randomNumber)}
		if output.job.id == 25 {
			cancel()
			return
		}
		results <- output
	}
}

func createWorkerPool(ctx context.Context, cancel context.CancelFunc, noOfWorkers int, insideJob chan Job, insideResult chan Result) {
	defer func() {
		close(insideResult)
		cancel()
	}()

	wg := sync.WaitGroup{}
	for i := 0; i < noOfWorkers; i++ {
		select {
		case <-ctx.Done():
			return
		default:
		}
		wg.Add(1)
		go workerReal(ctx, cancel, &wg, insideJob, insideResult)
	}
	wg.Wait()
}

func allocate(ctx context.Context, wg *sync.WaitGroup, jobs chan<- Job, noOfJobs int) {
	defer func() {
		close(jobs)
		wg.Done()
	}()

	for i := 0; i < noOfJobs; i++ {
		select {
		case <-ctx.Done():
			return
		default:
		}
		randomno := rand.Intn(999)
		job := Job{i, randomno}
		jobs <- job
	}
}

func result(ctx context.Context, wg *sync.WaitGroup, results <-chan Result, id string) {
	defer func() {
		wg.Done()
	}()

	for result := range results {
		select {
		case <-ctx.Done():
			return
		default:
		}
		log.Printf("%s: Job id %d, input random no %d , sum of digits %d\n", id, result.job.id, result.job.randomNumber, result.sumOfDigits)
	}
}

func main() {
	insideJob := make(chan Job, 10)

	defer func() {
		close(insideJob)
		log.Println("asdasda")
	}()

	return

	wg := sync.WaitGroup{}
	wg.Add(1)
	go realPoolingGoroutine("F1", &wg)
	wg.Add(1)
	go realPoolingGoroutine("F2", &wg)

	wg.Wait()
}

func realPoolingGoroutine(id string, wgs *sync.WaitGroup) {

	insideJob := make(chan Job, 10)
	insideResult := make(chan Result, 10)

	startTime := time.Now()
	noOfJobs := 100

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		wgs.Done()
	}()

	noOfWoker := 10
	go createWorkerPool(ctx, cancel, noOfWoker, insideJob, insideResult)

	wgx := sync.WaitGroup{}
	wgx.Add(1)
	go allocate(ctx, &wgx, insideJob, noOfJobs)
	wgx.Add(1)
	go result(ctx, &wgx, insideResult, id)

	wgx.Wait()

	diff := time.Since(startTime)
	log.Println("time taken", diff)
}

func semver01() {
	jobs := make(chan int, 100)
	results := make(chan int, 100)

	var wg sync.WaitGroup

	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go workerNew(w, &wg, jobs, results)
	}

	close(jobs)
	wg.Wait()
	return

	for j := 1; j <= 105; j++ {
		jobs <- j
	}

	close(jobs)

	for a := 1; a <= 5; a++ {
		log.Println(<-results)
	}

	wg.Wait()

	close(results)
}

func workerNew(id int, wg *sync.WaitGroup, jobs <-chan int, results chan<- int) {
	defer func() {
		log.Println("ended", id)
		wg.Done()
	}()

	log.Println("called worker:", id)

	for j := range jobs {
		log.Println("worker", id, "started job", j)
		time.Sleep(time.Second)
		log.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}

func nothing() {
	defer func() {
		log.Println("does nothing")
	}()
}

func worker(done chan bool) {
	log.Println("Working...")
	time.Sleep(time.Second)
	log.Println("done")

	done <- true
}

func syncChannel() {
	jobs := make(chan int, 5)
	done := make(chan bool)

	var wg sync.WaitGroup

	go func() {
		for {
			j, more := <-jobs
			log.Println(j, more)

			if more {
				log.Println("recieved", j)
				wg.Done()
			} else {
				log.Println("recieved all jobs")
				done <- true
				return
			}
		}
	}()

	for j := 1; j <= 3; j++ {
		wg.Add(1)
		jobs <- j
		log.Println("send", j)
	}

	close(jobs)

	wg.Wait()
	log.Println("sent all jobs")

	log.Println(<-done)
}

func mutex() {
	var (
		state    = make(map[int]int)
		mutex    = &sync.Mutex{}
		readOps  uint64
		writeOps uint64
	)

	for read := 0; read < 100; read++ {
		go func() {
			total := 0
			for {
				key := rand.Intn(5)
				mutex.Lock()
				total += state[key]
				mutex.Unlock()
				atomic.AddUint64(&readOps, 1)
			}
		}()
	}

	for write := 0; write < 10; write++ {
		go func() {
			for {
				key := rand.Intn(5)
				val := rand.Intn(100)
				mutex.Lock()
				state[key] = val
				mutex.Unlock()
				atomic.AddUint64(&writeOps, 1)
			}
		}()
	}

	time.Sleep(time.Second)

	mutex.Lock()

	readOpsFinal := atomic.LoadUint64(&readOps)
	log.Println("readOps:", readOpsFinal)
	writeOpsFinal := atomic.LoadUint64(&writeOps)
	log.Println("writeOps:", writeOpsFinal)
	log.Println("state:", state)
	log.Println("###################################")
	mutex.Unlock()
	time.Sleep(time.Millisecond)

	readOpsFinal = atomic.LoadUint64(&readOps)
	log.Println("readOps:", readOpsFinal)
	writeOpsFinal = atomic.LoadUint64(&writeOps)
	log.Println("writeOps:", writeOpsFinal)
	log.Println("state:", state)
}
