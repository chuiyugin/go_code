package main

import (
	"first/calc"
	"fmt"

	"github.com/q1mi/hello"
)

func main() {
	fmt.Println("hello world!")
	hello.SayHi()
	fmt.Println(calc.Add(1, 2))
	var x interface{} // 定义一个空接口
	x = "yes"
	switch v := x.(type) {
	case bool:
		fmt.Printf("猜对了, %v\n", v)
	default:
		fmt.Printf("猜错了, %v\n", v)
	}
}
