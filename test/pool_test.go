package test

import (
	"errors"
	"strconv"
	"sync"
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
		return nil, errors.New("no more")
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
			_ = manPool.releaseMan(man)
		} else {
			t.Log("no man")
		}
	}
}

func TestSyncPool(t *testing.T) {
	pool := &sync.Pool{
		New: func() interface{} {
			return 0
		},
	}
	t.Log(pool.Get().(int))
	pool.Put(20)
	t.Log(pool.Get().(int))
	t.Log(pool.Get().(int)) //0，取出则消失
	t.Log(pool)
}

func TestSyncPool2(t *testing.T) {
	pool := &sync.Pool{
		New: func() interface{} {
			return 10
		},
	}
	pool.Put(11) //私有对象，必定首次获取到
	pool.Put(12)
	pool.Put(13)
	pool.Put(14)
	pool.Put(15)
	for i := 0; i < 10; i++ {
		go func() {
			t.Log(pool.Get().(int))
		}()
	}
	select {}
}
