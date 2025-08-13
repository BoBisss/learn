package main

import (
	"fmt"
	"sync"
	"time"
)

var wg = sync.WaitGroup{}

func main() {
	//题目 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
	questionOne()

	//题目 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
	questionTwo()
}

func questionOne() {
	//打印偶数
	go func() {
		fmt.Print("偶数值有：")
		for i := 2; i <= 10; i++ {
			if i%2 == 0 {
				fmt.Printf("%d\t", i)
			}
		}
		fmt.Println()
	}()
	//打印奇数
	go func() {
		fmt.Print("奇数值有：")
		for i := 1; i <= 10; i++ {
			if i%2 != 0 {
				fmt.Printf("%d\t", i)
			}
		}
		fmt.Println()
	}()
	time.Sleep(2 * time.Second)
}

func questionTwo() {
	func1 := func() {
		time.Sleep(500 * time.Millisecond)
	}
	func2 := func() {
		time.Sleep(200 * time.Millisecond)
	}
	func3 := func() {
		time.Sleep(6000 * time.Millisecond)
	}
	func4 := func() {
		time.Sleep(789 * time.Millisecond)
	}
	func5 := func() {
		time.Sleep(999 * time.Millisecond)
	}
	func6 := func() {
		time.Sleep(7000 * time.Millisecond)
	}
	now := time.Now()
	tasks := map[int]func(){1: func1, 2: func2, 3: func3, 4: func4, 5: func5, 6: func6}
	times := task(tasks)
	for key, value := range times {
		fmt.Printf("任务%d所耗费的时间是：%v\n", key, value)
	}
	fmt.Println("总共耗时：", time.Since(now))

}

func task(tasks map[int]func()) (times map[int]time.Duration) {
	times = make(map[int]time.Duration)
	for key, fun := range tasks {
		wg.Add(1)
		go timer(times, key, fun)
	}
	wg.Wait()
	return
}

func timer(times map[int]time.Duration, key int, fun func()) {
	defer wg.Done()
	start := time.Now()
	fun()
	duration := time.Since(start)
	times[key] = duration
}
