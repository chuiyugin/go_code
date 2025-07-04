package main // 定义当前包名

// fmt包实现了格式化输出的函数
import (
	"fmt"
	"math"
)

func main() {
	// 一行代码代表语句结束
	fmt.Println("hello yugin!")
	// 标准声明格式
	var name string
	var age int
	var isOK bool
	fmt.Println(name, age, isOK)
	// 批量声明
	// var (
	// 	a string
	// 	b int
	// 	c bool
	// 	d float32
	//)
	//fmt.Println(a, b, c, d)
	// 常量声明
	const (
		n1 = iota
		n2
		_
		n3
	)
	fmt.Println(n1, n2, n3)
	const (
		a, b = iota + 1, iota + 2 //1,2
		c, d                      //2,3
		e, f                      //3,4
	)
	fmt.Println(a, b, c, d, e, f)
	fmt.Println(math.MaxFloat64)
}
