package main

import (
	"sync"
)

type MyChan struct {
	sync.Once
	done chan struct{}
}

func (m *MyChan) Close() {
	m.Do(func() {
		close(m.done)
	})
}

func NewMyChan() *MyChan {
	return &MyChan{
		done: make(chan struct{}),
	}
}

type ConcurrencyMap struct {
	sync.Mutex
	kv        map[int]int
	mapToChan map[int]*MyChan
}

func (m *ConcurrencyMap) Put(key, val int) {
	m.Lock()
	defer m.Unlock()
	m.kv[key] = val

	ch, ok := m.mapToChan[key]
	if !ok {
		return
	}

	ch.Close()
}

func (m *ConcurrencyMap) Get(key int) (val int, ok bool) {
	m.Lock()
	val, ok = m.kv[key]
	if ok {
		return
	}

	var ch *MyChan
	ch, ok = m.mapToChan[key]
	if !ok {
		ch = NewMyChan()
		m.mapToChan[key] = ch
	}

	m.Unlock()
	<-ch.done
	m.Lock()
	defer m.Unlock()
	val, ok = m.kv[key]
	return
}
