package main

import (
	"fmt"

	"io/ioutil"
	"os"
	"strings"
)

func main() {
	counts := make(map[string]int)
	for _, filename := range os.Args[1:] {
		data, err := ioutil.ReadFile(filename)

		if err != nil {
			fmt.Fprintf(os.Stderr, "dup2:%v\n", err)
			continue
		}
		//windows下换行符要写成\r\n
		fmt.Println(string(data))
		for _, line := range strings.Split(string(data), "\r\n") {

			counts[line] = counts[line] + 1

		}
	}

	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%s\t%d\n", line, n)
		}
	}

}
