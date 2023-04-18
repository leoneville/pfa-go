package main

import (
	"fmt"
	"sync"
)

func contador() {
	for i := 0; i < 100; i++ {
		fmt.Println(i)
	}
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		contador()
	}()
	wg.Wait()
}
