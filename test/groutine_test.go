package test

import (
	"context"
	"fmt"
	"sync"
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