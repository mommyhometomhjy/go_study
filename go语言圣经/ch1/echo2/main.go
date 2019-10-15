package main

import (
	"fmt"
	"os"
)

func main() {

	for i, arg := range os.Args[1:] {
		fmt.Printf("第%d个参数为%s\n", i+1, arg)
	}

}
