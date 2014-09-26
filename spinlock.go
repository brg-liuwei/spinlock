package spinlock

/*
警告：如果在单核环境下((runtime.NumCPU() == 1)
或者没有设置runtime.GOMAXPROCS，使用此自旋锁会带来性能下降
*/

/*
static void memory_barrier() {
    __sync_synchronize();
}
*/
import "C"

import (
	//"runtime"
	"sync/atomic"
)

type SpinLock struct {
	lock int32
}

func memoryBarrier() {
	//C.memory_barrier()
}

const MAXLOOP uint32 = 0x1 << 20

func (spin *SpinLock) TryLock() bool {
	return atomic.CompareAndSwapInt32(&spin.lock, 0, 1)
}

func (spin *SpinLock) Lock() {
	var i, loop uint32
	if spin.TryLock() {
		return
	}
	for {
		for loop = 1; loop != MAXLOOP; loop <<= 1 {
			for i = 0; i < loop; i++ {
				memoryBarrier()
			}
			if spin.TryLock() {
				return
			}
		}
		//runtime.Gosched()
	}
}

func (spin *SpinLock) UnLock() {
	if atomic.AddInt32(&spin.lock, -1) != 0 {
		panic("SpinLock Unlock invoke error")
	}
}

type SpinRWLock struct {
	r int32
	w int32
}

func (spin *SpinRWLock) TryWLock() bool {
	/* test rLock */
	if !atomic.CompareAndSwapInt32(&spin.r, 0, 0) {
		return false
	}
	/* try to acquire wLock */
	if !atomic.CompareAndSwapInt32(&spin.w, 0, 1) {
		return false
	}
	/* check rLock */
	if atomic.CompareAndSwapInt32(&spin.r, 0, 0) {
		return true
	} else {
		atomic.AddInt32(&spin.w, -1)
		return false
	}
}

func (spin *SpinRWLock) WLock() {
	var i, loop uint32
	if spin.TryWLock() {
		return
	}
	for {
		for loop = 1; loop != MAXLOOP; loop <<= 1 {
			for i = 0; i < loop; i++ {
				memoryBarrier()
			}
			if spin.TryWLock() {
				return
			}
		}
		//runtime.Gosched()
	}
}

func (spin *SpinRWLock) UnWLock() {
	if atomic.AddInt32(&spin.w, -1) != 0 {
		panic("SpinRWLock UnLock Error")
	}
}

func (spin *SpinRWLock) TryRLock() bool {
	if !atomic.CompareAndSwapInt32(&spin.w, 0, 0) {
		return false
	}
	atomic.AddInt32(&spin.r, 1)
	if !atomic.CompareAndSwapInt32(&spin.w, 0, 0) {
		return true
	} else {
		atomic.AddInt32(&spin.r, -1)
		return false
	}
}

func (spin *SpinRWLock) RLock() {
	var i, loop uint32
	if spin.TryRLock() {
		return
	}
	for {
		for loop = 1; loop != MAXLOOP; loop <<= 1 {
			for i = 0; i < loop; i++ {
				memoryBarrier()
			}
			if spin.TryRLock() {
				return
			}
		}
		//runtime.Gosched()
	}
}

func (spin *SpinRWLock) UnRLock() {
	if atomic.AddInt32(&spin.r, -1) < 0 {
		panic("SpinRWLock UnRLock Error")
	}
}
