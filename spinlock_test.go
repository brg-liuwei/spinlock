package spinlock

import (
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

var LOOPS int = 10

func TestMutexLock(t *testing.T) {
	// test mutex lock
	runtime.GOMAXPROCS(runtime.NumCPU())
	var wg0 sync.WaitGroup
	var lock sync.Mutex
	v := 0

	now := time.Now()
	for i := 0; i < LOOPS; i++ {
		wg0.Add(1)
		go func(val *int, lock *sync.Mutex, seq int) {
			{
				lock.Lock()
				defer lock.Unlock()
				if *val != 0 {
					t.Error("*val = ", *val)
				}
				for j := 0; j != 10; j++ {
					*val += rand.Int()
				}
				*val = 0
			}
			wg0.Done()
		}(&v, &lock, i)
	}
	wg0.Wait()
	println("MutexLock use time: ", time.Since(now))
}

func TestSpinLock(t *testing.T) {
	if runtime.NumCPU() == 1 {
		println("runtime.NumCPU() == 1, Unsuggest using this package")
		return
	} else {
		println("cpu: ", runtime.NumCPU())
		println("set proc: ", runtime.GOMAXPROCS(runtime.NumCPU()))
	}

	var data int
	var spin SpinLock
	var wg sync.WaitGroup

	now := time.Now()
	for i := 0; i < LOOPS; i++ {
		wg.Add(1)
		go func(val *int, lock *SpinLock, seq int) {
			{
				lock.Lock()
				defer lock.UnLock()
				if *val != 0 {
					t.Error("*val = ", *val)
				}
				for j := 0; j != 10; j++ {
					*val += rand.Int()
				}
				*val = 0
			}
			wg.Done()
		}(&data, &spin, i)
	}
	wg.Wait()
	println("SpinLock use time: ", time.Since(now))
}
