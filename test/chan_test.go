package test

import (
	"fmt"
	"os"
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
	resChan := writeFieldToFile()
	resChan <- []string{"aa", "bb"}
	t.Log("aabb")
	resChan <- []string{"cc", "dd"}
	t.Log("aabb")
	close(resChan)
}

func writeFieldToFile() chan []string {
	out := make(chan []string)
	go func() {
		if dbFile, err := os.OpenFile("tmp", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755); err == nil {
			for item := range out {
				fmt.Println("接受：", item)
				for _, str := range item {
					if _, err := dbFile.Write([]byte(str)); err != nil {
						fmt.Println("database写入失败" + err.Error())
					}
				}
				time.Sleep(3 * time.Second)
			}
			dbFile.Close()
		} else {
			panic("fields.csv打开失败：" + err.Error())
		}
	}()
	return out
}
