package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type ChanMutex chan struct{}

func NewTryLock() ChanMutex {
	ch := make(chan struct{}, 1)
	return ch
}
func (m *ChanMutex) Lock() {
	ch := (chan struct{})(*m)
	ch <- struct{}{}
}
func (m *ChanMutex) Unlock() {
	ch := (chan struct{})(*m)
	select {
	case <-ch:
	default:
		panic("unlock of unlocked mutex")
	}
}
func (m *ChanMutex) TryLockWithTimeOut(d time.Duration) bool {
	ch := (chan struct{})(*m)
	t := time.NewTimer(d)
	select {
	case <-t.C:
		return false
	case ch <- struct{}{}:
		t.Stop()
		return true
	}
}
func (m *ChanMutex) TryLock() bool {
	ch := (chan struct{})(*m)
	select {
	case ch <- struct{}{}:
		return true
	default:
		return false
	}
}

func main() {
	n1 := int64(0)
	n2 := int64(0)
	c := NewTryLock()

	wg := sync.WaitGroup{}
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			if c.TryLock() {
				n1++
				c.Unlock()
			} else {
				atomic.AddInt64(&n2, 1)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	fmt.Printf("total: %v, success: %v, fail: %v\n", n1+n2, n1, n2)
}
