package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"
)

func TestChan(t *testing.T) {
	resChan := make(chan int)
	go func() {
		for x := range resChan {
			fmt.Println(x)
		}
	}()
	resChan <- 1
	resChan <- 2
	resChan <- 3
	//close(resChan)
	//return;
}

func TestChanWrite(t *testing.T) {
	var wgWrite sync.WaitGroup
	wgWrite.Add(1)
	resChan := WriteFieldToFile(&wgWrite)
	resChan <- []string{"aa", "bb"}
	t.Log("aabb")
	resChan <- []string{"cc", "dd"}
	t.Log("aabb")
	close(resChan)
	wgWrite.Wait()
	if content, err := ioutil.ReadFile("tmp"); err == nil {
		t.Log(string(content))
		_ = os.Remove("tmp")
	} else {
		t.Log(err.Error())
	}
}

func WriteFieldToFile(wg *sync.WaitGroup) chan []string {
	out := make(chan []string)
	go func() {
		defer wg.Done()
		if dbFile, err := os.OpenFile("tmp", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755); err == nil {
			for item := range out {
				fmt.Println("接受：", item)
				for _, str := range item {
					if _, err := dbFile.Write([]byte(str)); err != nil {
						fmt.Println("database写入失败" + err.Error())
					}
				}
			}
			_ = dbFile.Close()
		} else {
			panic("fields.csv打开失败：" + err.Error())
		}
	}()
	return out
}

func TestSelectChan(t *testing.T) {
	chan1 := make(chan string)
	go func() {
		fmt.Println("sleep1")
		time.Sleep(time.Second * 2)
		//chan1 <- "bbbb"
		fmt.Println("sleep2")
	}()
	select { // 阻塞选择，选择完才继续走下去
	case ret := <-chan1:
		t.Log(ret)
	case <-time.After(time.Second):
		fmt.Println("timeout")
	}
	fmt.Println("select结束")
	//无缓冲通道，所以当对这个缓冲通道写的时候，
	//会一直阻塞等到某个协程对这个缓冲通道读（大家发现没有这个与典型的生产者消费者有点不一样，当队列中“内容”已经满了，
	//生产者再生往里放东西才会阻塞，而这里我讲c<-'A’理解为生产，他却是需要等到某个协程读了再能继续运行）。
	//chan1 <- "aaaa"//all goroutines are asleep - deadlock
}

func TestSelectChanClose(t *testing.T) {
	chan1 := make(chan string)
	fmt.Println(len(chan1)) // len为0，阻塞式通道
	go func() {
		fmt.Println("sleep1")
		time.Sleep(time.Second * 1)
		close(chan1)
		fmt.Println("sleep2")
	}()
	select {
	case ret, ok := <-chan1:
		t.Log(ret, ok)
	case <-time.After(time.Second * 2):
		fmt.Println("timeout")
	}
	fmt.Println("select结束")
}

func TestBufferedChan(t *testing.T) {
	chan1 := make(chan int, 5)
	chan1 <- 1
	chan1 <- 2
	//chan1 <- 5//all goroutines are asleep - deadlock
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			x, ok := <-chan1
			fmt.Println(x, ok)
			time.Sleep(1 * time.Second)
		}
	}()
	chan1 <- 3
	chan1 <- 4
	chan1 <- 5
	chan1 <- 6 // 满了，阻塞
	chan1 <- 7
	fmt.Println("wait")
	close(chan1)
	wg.Wait()
}
