package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	/*
		题目 ：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。
		启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
		考察点 ： sync.Mutex 的使用、并发数据安全。
	*/
	questionOne()
	/*
		题目 ：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。
		启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
		考察点 ：原子操作、并发数据安全。
	*/
	questionTwo()

}

var lock = sync.Mutex{}

func questionOne() {
	var wg = sync.WaitGroup{}
	num := 0
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				increase(&num)
			}
		}()
	}
	wg.Wait()
	fmt.Println("结果是：", getNum(&num))
}

func increase(num *int) {
	lock.Lock()
	*num++
	lock.Unlock()
}
func getNum(num *int) int {
	lock.Lock()
	defer lock.Unlock()
	return *num
}

func questionTwo() {
	var wg = sync.WaitGroup{}
	num := int64(0)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				atomic.AddInt64(&num, 1)
			}
		}()
	}
	wg.Wait()
	fmt.Println("结果是：", atomic.LoadInt64(&num))
}
