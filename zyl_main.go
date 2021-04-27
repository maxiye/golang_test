package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main1() {
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

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}

func main() {
	http.HandleFunc("/", IndexHandler)
	http.ListenAndServe("127.0.0.1:8988", nil)
}
