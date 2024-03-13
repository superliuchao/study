package main

import (
	"bytes"
	"encoding/binary"
	"sync"
)

type part struct {
	sync.Mutex
	mp map[int]int
}

type Pmap struct {
	p      []*part
	pcount int
}

func IntToBytes(n int) []byte {
	data := uint64(n)
	bytebuf := &bytes.Buffer{}
	binary.Write(bytebuf, binary.BigEndian, data)
	return bytebuf.Bytes()
}

func BytesToInt(bys []byte) int {
	bytebuff := bytes.NewBuffer(bys)
	var data int
	binary.Read(bytebuff, binary.BigEndian, &data)
	return data
}

func New(partitionCount int) *Pmap {
	m := &Pmap{
		p:      make([]*part, partitionCount),
		pcount: partitionCount,
	}

	for i := 0; i < partitionCount; i++ {
		m.p[i] = &part{
			mp: make(map[int]int),
		}
	}
	return m
}

func (m *Pmap) Put(key, val int) {
	//mur32 := murmur3.New32().Sum(IntToBytes(uint64(key)))
	//BytesToInt()
	index := uint(key) % uint(m.pcount)
	p := m.get(int(index))
	p.Lock()
	m.p[index].mp[key] = val
	p.Unlock()
	return
}

func (m *Pmap) Get(key int) (val int, ok bool) {

	index := int(uint(key) % uint(m.pcount))
	p := m.get(index)
	p.Lock()
	defer p.Unlock()
	val, ok = p.mp[key]
	return
}

func (m *Pmap) get(index int) *part {
	return m.p[index]
}
