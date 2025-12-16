package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ExpireLock 带过期时间的锁结构体
type ExpireLock struct {
	// 核心锁
	mutex sync.Mutex
	// 流程锁，防止多次解锁，例如异步解锁协程解锁和手动解锁同时发生
	processMutex sync.Mutex
	// 锁的身份标识
	token string
	// 停止异步解锁协程函数
	stop context.CancelFunc
}

// NewExpireLock 创建新的过期锁实例
func NewExpireLock() *ExpireLock {
	return &ExpireLock{}
}

// Lock 加锁
func (e *ExpireLock) Lock(expireTime int64) {
	// 1. 加锁
	e.mutex.Lock()

	// 2. 设置锁的身份标识token
	token := GetProcessAndGoroutineIDStr()
	e.token = token

	// 2.1 校验过期时间，如果小于等于0，代表手动释放锁，无需开启异步解锁协程
	if expireTime <= 0 {
		return
	}

	// 3. 给终止异步协程函数stop赋值，启动异步协程，达到指定时间后执行解锁操作
	ctx, cancel := context.WithCancel(context.Background())
	e.stop = cancel

	go func() {
		select {
		// 到了锁的过期时间，释放锁
		case <-time.After(time.Duration(expireTime) * time.Second):
			e.unlock(token)
		case <-ctx.Done():
		}
	}()
}

// Unlock 解锁
func (e *ExpireLock) Unlock() error {
	return e.unlock(GetProcessAndGoroutineIDStr())
}

// unlock 内部解锁方法
func (e *ExpireLock) unlock(token string) error {
	// 1. 加流程锁，防止并发情况下，异步解锁协程解锁和手动解锁同时发生
	e.processMutex.Lock()
	defer e.processMutex.Unlock()

	// 2. 校验token
	if e.token != token {
		return errors.New("unlock not your lock")
	}

	// 3. 停止异步解锁协程
	if e.stop != nil {
		e.stop()
	}

	// 4. 重置token
	e.token = ""

	// 5. 解锁
	e.mutex.Unlock()

	return nil
}

// GetCurrentProcessID 获取当前进程ID
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

// GetProcessAndGoroutineIDStr 获取进程ID和协程ID的组合字符串
func GetProcessAndGoroutineIDStr() string {
	return fmt.Sprintf("%s_%s", GetCurrentProcessID(), GetCurrentGoroutineID())
}

func main() {
	// 测试带过期时间的锁
	lock := NewExpireLock()

	fmt.Println("测试带过期时间的锁...")

	// 加锁，设置5秒后自动过期
	lock.Lock(5)
	fmt.Println("锁已获取，将在5秒后自动释放")

	// 模拟一些工作
	time.Sleep(2 * time.Second)
	fmt.Println("工作完成，手动释放锁")

	// 手动释放锁
	err := lock.Unlock()
	if err != nil {
		fmt.Printf("解锁失败: %v\n", err)
	} else {
		fmt.Println("锁已成功释放")
	}

	// 测试手动释放锁（过期时间为0）
	fmt.Println("\n测试手动释放锁...")
	lock2 := NewExpireLock()
	lock2.Lock(0) // 过期时间为0，需要手动释放
	fmt.Println("锁已获取，需要手动释放")

	time.Sleep(1 * time.Second)
	err = lock2.Unlock()
	if err != nil {
		fmt.Printf("解锁失败: %v\n", err)
	} else {
		fmt.Println("锁已成功释放")
	}
}
