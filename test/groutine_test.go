package test

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestGo(t *testing.T) {
	for i := 0; i < 10; i++ {
		/*go func() {
			t.Log(i)//乱
		}()*/
		go func(ii int) {
			t.Log(ii) //正
		}(i)
	}
	time.Sleep(time.Millisecond * 500)
}

func TestRwLock(t *testing.T) {
	var rwLock sync.RWMutex
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		rwLock.RLock()
		t.Log("RLock", time.Now().Unix())
		time.Sleep(time.Second * 2)
		rwLock.RUnlock()
		wg.Done()
	}()
	go func() {
		rwLock.RLock()
		t.Log("RLock", time.Now().Unix())
		time.Sleep(time.Second * 2)
		rwLock.RUnlock()
		wg.Done()
	}()
	go func() {
		rwLock.Lock()
		t.Log("Lock", time.Now().Unix())
		time.Sleep(time.Second * 2)
		rwLock.Unlock()
		wg.Done()
	}()
	wg.Wait()
}

func TestGoOpt(t *testing.T) {
	counter := 0
	var mu sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < 5000; i++ {
		wg.Add(1)
		go func() {
			defer func() {
				mu.Unlock()
				wg.Done()
			}()
			mu.Lock()
			counter++ // 不加锁结果异常
		}()
	}
	//time.Sleep(time.Millisecond * 500)
	wg.Wait()
	t.Log(counter)
}

func TestAtomicAdd(t *testing.T) {
	var num int64
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			atomic.AddInt64(&num, 1)
			wg.Done()
		}()
	}
	wg.Wait()
	t.Log(num)
}

func TestTimer(t *testing.T) {
	timer1 := time.NewTimer(5 * time.Second)
	ticker1 := time.NewTicker(2 * time.Second)
	time.Sleep(6 * time.Second)
	fmt.Println("sleep done")
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			t, err := <-timer1.C
			fmt.Println("timer do: ", t.Format("2006/01/02 15:04:05"), err)
			timer1.Reset(1 * time.Second)
		}
	}()

	go func() {
		defer wg.Done()
		for {
			t, err := <-ticker1.C
			fmt.Println("ticker do: ", t.Format("2006/01/02 15:04:05"), err)
		}
	}()
	wg.Wait()
}

func TestCtx(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background()) // 不可变对象
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Done")
				return
			case <-time.After(1 * time.Second):
				fmt.Println("1s later...")
			}
		}
	}()
	time.Sleep(5 * time.Second)
	cancel()
	wg.Wait()
}

func TestAsyncTask(t *testing.T) {
	asyncTask := func() chan string {
		resChan := make(chan string)
		t.Log(len(resChan))
		go func() {
			time.Sleep(time.Second * 2)
			resChan <- "result"
		}()
		return resChan
	}
	resChan := asyncTask()
	t.Log(<-resChan)
}

func TestStatefull(t *testing.T) {
	type readOp struct {
		key  int
		resp chan int
	}
	type writeOp struct {
		key  int
		val  int
		resp chan bool
	}

	var readOps uint64
	var writeOps uint64

	reads := make(chan readOp)
	writes := make(chan writeOp)

	go func() {
		var state = make(map[int]int)
		for {
			select {
			case read := <-reads:
				read.resp <- state[read.key]
			case write := <-writes:
				state[write.key] = write.val
				write.resp <- true
			}
		}
	}()

	for r := 0; r < 100; r++ {
		go func() {
			for {
				read := readOp{
					key:  rand.Intn(5),
					resp: make(chan int)}
				reads <- read
				<-read.resp
				atomic.AddUint64(&readOps, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	for w := 0; w < 10; w++ {
		go func() {
			for {
				write := writeOp{
					key:  rand.Intn(5),
					val:  rand.Intn(100),
					resp: make(chan bool)}
				writes <- write
				<-write.resp
				atomic.AddUint64(&writeOps, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	time.Sleep(time.Second)

	readOpsFinal := atomic.LoadUint64(&readOps)
	fmt.Println("readOps:", readOpsFinal)
	writeOpsFinal := atomic.LoadUint64(&writeOps)
	fmt.Println("writeOps:", writeOpsFinal)
}

func TestLock(t *testing.T) {
	var lock sync.Mutex
	go func() {
		lock.Lock()
	}()
}
