package main

import "fmt"

func main() {
	//写入defer关键字
	defer func() {
		//查看是否有异常
		if errs := recover(); errs != nil {
			fmt.Printf("main end3, err= %v", errs)
		}
	}()
	defer func() {
		fmt.Println("main end1")
	}()
	defer func() {
		fmt.Println("main end2")
	}()
	panic("error panic")
}
