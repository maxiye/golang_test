package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestChan(t *testing.T) {
	resChan := make(chan int64) //阻塞式chan，上一个未pop，无法push下一个
	go func() {
		for x := range resChan {
			fmt.Println(x)
			time.Sleep(time.Second)
		}
	}()
	resChan <- time.Now().Unix()
	resChan <- time.Now().Unix()
	resChan <- time.Now().Unix()
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
	chan1 := make(chan int, 3)
	chan1 <- 1
	chan1 <- 2
	chan1 <- 3
	//chan1 <- 5//all goroutines are asleep - deadlock
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			x, ok := <-chan1
			if !ok {
				fmt.Println("chan closed")
				break
			}
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

func TestSelectHold(t *testing.T) {
	go func() {
		ticker := time.NewTicker(time.Second * 1)
		wait := 5
		for x := range ticker.C {
			t.Log(x.Unix())
			wait--
			if wait == 0 {
				//break//仍然阻塞中，需要手动停
				os.Exit(0)
			}
		}
	}()
	select {} //block 进程，没有协程在后台运行就panic：fatal error: all goroutines are asleep - deadlock!
}

func TestSelectDefault(t *testing.T) {
	//chan1 := make(chan int)//fatal error: all goroutines are asleep - deadlock!
	chan1 := make(chan int, 2)
	chan1 <- 1
	chan2 := make(chan string, 1)
	chan2 <- "a"
	times := 5
	for {
		select { // 随机一个chan获取，无上下顺序
		case res := <-chan1:
			t.Log("chan1", res)
		case res := <-chan2:
			t.Log("chan2", res)
		default: //非阻塞，不等待
			t.Log("default")
			time.Sleep(time.Millisecond * 500)
			times--
			if times == 0 {
				//break
				//continue// break、continue进入下一次select
				goto out
			}
		}
	}
out:
}

// 通道方向
func TestChanDirection(t *testing.T) {
	ping := func(pings chan<- string, msg string) {
		//msg2 := <- pings// invalid operation: <-pings (receive from send-only type chan<- string)
		pings <- msg
	}
	pong := func(pings <-chan string, pongs chan<- string) {
		msg := <-pings
		pongs <- msg
	}
	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	go func() {
		for {
			pong(pings, pongs)
			fmt.Println(<-pongs) // 必须释放pongs，不然阻塞
		}
	}()
	ping(pings, "1")
	ping(pings, "2")
	ping(pings, "3")
	time.Sleep(time.Second * 1)
}

// 通道遍历
func TestRangeChan(t *testing.T) {
	chan1 := make(chan int, 2)
	chan1 <- 1
	chan1 <- 2
	close(chan1)
	// 一个非空的通道也是可以关闭的，并且，通道中剩下的值仍然可以被接收到。
	for v := range chan1 {
		t.Log(v)
	}
}

// pool
func TestWorker(t *testing.T) {
	worker := func(id int, jobs <-chan string, res chan<- string) {
		for job := range jobs {
			t.Logf("worker:%d do job:%s", id, job)
			time.Sleep(time.Second * 2)
			res <- job + ":" + strconv.FormatInt(time.Now().Unix(), 10)
			t.Logf("worker:%d over job:%s", id, job)
		}
	}
	jobs := make(chan string, 5)
	res := make(chan string, 5)
	for i := 0; i < 3; i++ {
		go worker(i, jobs, res)
	}
	for i := 0; i < 5; i++ {
		jobs <- fmt.Sprintf("task-%d", i)
	}
	close(jobs)          // 结束通道，同时worker停止
	for v := range res { // fatal error: all goroutines are asleep - deadlock!
		t.Log(v)
	}
}

func TestNilChan(t *testing.T) {
	var strchan chan string
	strchan <- "aa"
}

var testChan chan string

func GetTestChan() chan string {
	return testChan
}

func TestGetChan(t *testing.T) {
	chan1 := GetTestChan()
	t.Log(chan1)
}
