package test

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestMutex(t *testing.T) {
	var lock sync.Mutex
	go func() {
		lock.Lock()
		fmt.Println("sleep 3")
		time.Sleep(3 * time.Second)
		lock.Unlock()
	}()
	go func() {
		lock.Lock()
		fmt.Println("sleep 2")
		time.Sleep(2 * time.Second)
		lock.Unlock()
	}()
	time.Sleep(6 * time.Second)
}
