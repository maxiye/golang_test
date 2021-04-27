package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main3() {
	uidVerFile := "I:\\test\\logs\\版本活跃用户810.csv"
	uidSet := make(map[string]bool, 150000)
	if uidVer, err := os.OpenFile(uidVerFile, os.O_RDONLY, os.ModePerm); err == nil {
		br := bufio.NewReader(uidVer)
		_, _ = br.ReadBytes('\n')
		for {
			l, e := br.ReadBytes('\n')
			if e != nil && len(l) == 0 {
				break
			}
			line := string(l)
			lineData := strings.Split(line, ",")
			uid, ver := lineData[0], lineData[1]
			if ver == "6.3" || ver == "6.3.1" || ver == "6.3.2" {
				uidSet[uid] = true
			}
		}
		fmt.Println(len(uidSet))
	}
	rfmFilePath := "I:\\test\\logs\\分层模型用户数据报表_810.csv"
	if rfm, err := os.OpenFile(rfmFilePath, os.O_RDONLY, os.ModePerm); err == nil {
		br := bufio.NewReader(rfm)
		_, _ = br.ReadBytes('\n')
		tMapCount := make(map[string]int)
		for {

			l, e := br.ReadBytes('\n')

			if e != nil && len(l) == 0 {
				break
			}
			line := string(l)
			lineData := strings.Split(line, ",")
			uid, t := lineData[2], lineData[9]
			if uidSet[uid] {
				tMapCount[t]++
			}
		}
		fmt.Println(tMapCount)
	}
}
