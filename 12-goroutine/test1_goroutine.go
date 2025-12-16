package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 主goroutine
func main() {
	var wg sync.WaitGroup

	// 打印初始状态
	fmt.Printf("程序启动时：\n")
	fmt.Printf("  Goroutines数量: %d\n", runtime.NumGoroutine())
	fmt.Printf("  CPU核心数: %d\n", runtime.NumCPU())

	fmt.Printf("主goroutine ID: %s\n", GetCurrentGoroutineID())

	// 创建多个goroutine来观察ID分配模式
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// 在子goroutine内部检查数量
			fmt.Printf("子goroutine %d: %s (当前goroutines: %d)\n", id, GetProcessAndGoroutineIDStr(), runtime.NumGoroutine())
		}(i)
	}

	// 立即检查goroutine数量
	fmt.Printf("\n创建goroutine后立即检查：\n")
	fmt.Printf("  Goroutines数量: %d\n", runtime.NumGoroutine())

	// 等待一下让goroutine有机会启动
	time.Sleep(50 * time.Millisecond)

	// 再次检查
	fmt.Printf("\n等待50ms后：\n")
	fmt.Printf("  Goroutines数量: %d\n", runtime.NumGoroutine())

	wg.Wait() // 等待所有goroutine完成

	// 打印最终状态
	fmt.Printf("\n程序结束时：\n")
	fmt.Printf("  Goroutines数量: %d\n", runtime.NumGoroutine())
}

func GetCurrentProcessID() string {
	return strconv.Itoa(os.Getpid())
}

// GetCurrentGoroutineID 获取当前的协程ID
func GetCurrentGoroutineID() string {
	buf := make([]byte, 128)
	buf = buf[:runtime.Stack(buf, false)]
	stackInfo := string(buf)
	return strings.TrimSpace(strings.Split(strings.Split(stackInfo, "[running]")[0], "goroutine")[1])
}

func GetProcessAndGoroutineIDStr() string {
	return fmt.Sprintf("%s_%s", GetCurrentProcessID(), GetCurrentGoroutineID())
}
