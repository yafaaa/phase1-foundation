package main

import (
	"fmt"
	"sync"
)

func main(){
	counter := 0
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i:=1; i <= 100; i++{
		wg.Add(1)
		
		go func(){
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()
			counter++

		}()
	}

	wg.Wait()
	fmt.Println("finished")

}