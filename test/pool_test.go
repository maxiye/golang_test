package test

import (
	"errors"
	"strconv"
	"testing"
	"time"
)

type ManPool struct {
	buffer chan *Man
	_      struct{}
}

func initPool(size int) *ManPool {
	pool := new(ManPool)
	pool.buffer = make(chan *Man, size)
	for i := 0; i < size; i++ {
		pool.buffer <- &Man{Name: "a" + strconv.Itoa(i), Age: byte(20 + i)}
	}
	return pool
}

func (pool *ManPool) getMan() (*Man, error) {
	select {
	case ret := <-pool.buffer:
		return ret, nil
	case <-time.After(1 * time.Second):
		return nil, errors.New("No more")
	}
}

func (pool *ManPool) releaseMan(man *Man) error {
	select {
	case pool.buffer <- man:
		return nil
	case <-time.After(time.Second):
		return errors.New("error release")
	}
}

func TestManPool(t *testing.T) {
	manPool := initPool(5)
	for i := 0; i < 6; i++ {
		man, err := manPool.getMan()
		if err == nil {
			man.Sleep(Girl{"aa", 20})
			manPool.releaseMan(man)
		} else {
			t.Log("no man")
		}
	}
}
