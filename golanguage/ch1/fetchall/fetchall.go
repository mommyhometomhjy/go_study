package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	start := time.Now()
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		go getBody(url, ch)
	}

	for range os.Args[1:] {
		fmt.Println(<-ch)
	}
	fmt.Printf("共耗时%.2fs", time.Since(start).Seconds())
}

func getBody(url string, ch chan<- string) {
	start := time.Now()
	if !strings.HasPrefix(url, "http") {
		url = "http" + url
	}
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	_, err = io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ch <- fmt.Sprintf("%s耗时%.2fs", url, time.Since(start).Seconds())
}
