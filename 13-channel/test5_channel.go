package main

import "fmt"

func fibonacii(c chan int, quit chan struct{}) {
	x, y := 1, 1

	for {
		select {
		case c <- x:
			//如果c可写，则该case就会进来
			x = y
			y = x + y
		case value := <-quit:
			fmt.Println("quit", value)
			return
		default:
			fmt.Println("default")
		}
	}
}

func main() {
	c := make(chan int)
	quit := make(chan struct{})

	//sub go
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c)
		}

		close(quit)
	}()

	//main go
	fibonacii(c, quit)
}
