package main

import (
	"context"
	"runtime"
	"strings"
	"sync"
	"time"
)

type expiredlock struct {
	mutex        sync.Mutex
	unlockCh     chan struct{}
	processMutex sync.Mutex
	id           string
	stop         context.CancelFunc
}

func GetRoutineID() string {
	//TODO 获取当前gorouteid
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	return idField
}

func (lock *expiredlock) Lock(expiredTimeout time.Duration) {

	lock.processMutex.Lock()
	defer lock.processMutex.Unlock()

	lock.mutex.Lock()
	lock.id = GetRoutineID()
	go func(id string) {
		ctx, stop := context.WithCancel(context.Background())
		lock.stop = stop
		select {
		case <-ctx.Done():

		case <-time.After(expiredTimeout):
			lock.unlock(id)
		}
	}(lock.id)

}

func (lock *expiredlock) unlock(id string) {
	lock.processMutex.Lock()
	defer lock.processMutex.Unlock()
	if lock.id == "" {
		return
	}

	if lock.id != id {
		return
	}
	lock.stop()
	lock.stop = nil
	lock.id = ""
	lock.mutex.Unlock()
	return
}

func (lock *expiredlock) Unlock() {
	lock.Unlock()
}
