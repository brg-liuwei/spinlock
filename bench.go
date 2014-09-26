package spinlock

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

func main() {
	loops := 0
	if len(os.Args) == 1 {
		loops = 100
	} else {
		var err error
		if loops, err = strconv.Atoi(os.Args[1]); err != nil {
			panic(err)
		}
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	now := time.Now()
	testMutexLock(loops)
	fmt.Println("test mutex lock: ", time.Since(now))
	now = time.Now()
	testSpinLock(loops)
	fmt.Println("test spin  lock: ", time.Since(now))
}

func testSpinLock(loop int) {
	var data int
	var spin SpinLock
	var wg sync.WaitGroup

	for i := 0; i < loop; i++ {
		wg.Add(1)
		go func(val *int, lock *SpinLock) {
			{
				lock.Lock()
				defer lock.UnLock()
				for j := 0; j != 10; j++ {
					*val += rand.Int()
				}
				*val = 0
			}
			wg.Done()
		}(&data, &spin)
	}
	wg.Wait()
}

func testMutexLock(loop int) {
	var wg sync.WaitGroup
	var lock sync.Mutex
	v := 0

	for i := 0; i < loop; i++ {
		wg.Add(1)
		go func(val *int, lock *sync.Mutex) {
			{
				lock.Lock()
				defer lock.Unlock()
				for j := 0; j != 10; j++ {
					*val += rand.Int()
				}
				*val = 0
			}
			wg.Done()
		}(&v, &lock)
	}
	wg.Wait()
}
