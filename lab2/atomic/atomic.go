package main

//声明这个文件属于 main包

import ( //导入程序所需的包
	"fmt"
	"sync"
	"sync/atomic"
)

// Global variables shared between functions --A BAD IDEA
var wg sync.WaitGroup //声明一个全局的 WaitGroup 变量

// 定义一个名为 addsAtomic的函数
func addsAtomic(n int, total *atomic.Int64) bool { //返回 bool类型：总是返回 true
	//n int：要执行的加法次数
	// total *atomic.Int64：指向一个原子整数的指针，用于安全地累加值
	for i := 0; i < n; i++ {
		total.Add(1)
	} //循环 n 次，每次将 total 的值原子地增加 1
	wg.Done() //let waitgroup know we have finished
	//通知 WaitGroup 这个 goroutine 已经完成
	return true
}

func main() {

	var total atomic.Int64
	//声明一个原子整数变量

	//for loop using range option
	for i := range 10 { //循环 10 次
		//the waitgroup is used as a barrier
		// init it to number of go routines
		wg.Add(1)
		fmt.Println("go Routine ", i)
		go addsAtomic(1000, &total)
	}
	wg.Wait() //wait here until everyone (10 go routines) is done
	fmt.Println(total.Load())

}
