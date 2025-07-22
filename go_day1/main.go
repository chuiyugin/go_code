package main // 定义当前包名

// fmt包实现了格式化输出的函数
import (
	"fmt"
)

func intSum2(x ...int) int {
	fmt.Println(x) //x是一个切片
	sum := 0
	for _, v := range x {
		sum = sum + v
	}
	return sum
}

// 定义全局变量num
var num int64 = 10

func testNum() {
	num := 1000
	fmt.Printf("num=%d\n", num) // 函数中优先使用局部变量
}

func main() {
	/*
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
		ret := intSum2(10, 20, 30, 40)
		fmt.Println(ret)
		testNum()
		fmt.Println(num)
	*/
	var point *int
	point = new(int)
	*point = 10
	fmt.Println(*point)
}
