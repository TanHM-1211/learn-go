package main

import "fmt"

func main() {
	done := make(chan int, 10)

	for i := 0; i < cap(done); i++ {
		go func(j int) {
			fmt.Println("Hello World", j)
			done <- 1
		}(i)
	}

	for i := 0; i < cap(done); i++ {
		<-done
	}
}
