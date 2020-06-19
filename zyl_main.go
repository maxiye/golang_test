package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("hello" + os.Args[1])
	} else {
		if t, err := time.Parse("2006-01-02", os.Args[1]); err == nil {
			fmt.Println("Second: ", t.Unix())
		} else {
			fmt.Println("ERROR, Now" + time.Now().Format("2006-01-02 15:04:05"))
		}
	}
	os.Exit(1)
}
