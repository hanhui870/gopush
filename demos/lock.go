package main

import (
	"sync"
	"fmt"
)

// error: locked 1 fatal error: all goroutines are asleep - deadlock!
// need reentrant lock
func main() {
	var lock sync.Mutex

	lock.Lock()
	fmt.Println("locked 1")

	lock.Lock()
	fmt.Println("locked 2")

	lock.Unlock()
	fmt.Println("unlocked 2")

	lock.Unlock()
	fmt.Println("unlocked 1")
}