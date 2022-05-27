// Copyright 2022 Guan Jianchang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yx

import "sync"

const (
	BP_MIN_BUFF_SIZE = 64
	BP_MAX_BUFF_SIZE = 4 * 1024
	BP_QUEUE_STEP    = 64
)

type BuffPool struct {
	pool     *sync.Pool
	buffSize uint32
}

func NewBuffPool(buffSize uint32) *BuffPool {
	p := &BuffPool{
		pool:     nil,
		buffSize: buffSize,
	}

	p.pool = &sync.Pool{
		New: p.newBuff,
	}

	return p
}

func (p *BuffPool) Get() *[]byte {
	return p.pool.Get().(*[]byte)
}

func (p *BuffPool) Put(buffRef *[]byte) {
	p.pool.Put(buffRef)
}

func (p *BuffPool) newBuff() interface{} {
	buff := make([]byte, p.buffSize)
	return &buff
}

type BuffFactory struct {
	poolList []*BuffPool
}

func NewBuffFactory() *BuffFactory {
	listCnt := (BP_MAX_BUFF_SIZE-BP_MIN_BUFF_SIZE)/BP_QUEUE_STEP + 1
	p := &BuffFactory{
		poolList: make([]*BuffPool, listCnt),
	}

	for i := 0; i < listCnt; i++ {
		size := uint32(BP_MIN_BUFF_SIZE + i*BP_QUEUE_STEP)
		p.poolList[i] = NewBuffPool(size)
	}

	return p
}

func (p *BuffFactory) CreateBuff(size uint32) *[]byte {
	if size > BP_MAX_BUFF_SIZE {
		buff := make([]byte, size)
		return &buff
	}

	idx := 0
	if size > BP_MIN_BUFF_SIZE {
		idx = (int(size) - BP_MIN_BUFF_SIZE) / BP_QUEUE_STEP
		if (size-BP_MIN_BUFF_SIZE)%BP_QUEUE_STEP != 0 {
			idx++
		}
	}

	buffRef := p.poolList[idx].Get()
	return buffRef
}

func (p *BuffFactory) ReuseBuff(buffRef *[]byte) {
	size := len(*buffRef)
	if size < BP_MIN_BUFF_SIZE || size > BP_MAX_BUFF_SIZE {
		return
	}

	if (size-BP_MIN_BUFF_SIZE)%BP_QUEUE_STEP != 0 {
		return
	}

	idx := (size - BP_MIN_BUFF_SIZE) / BP_QUEUE_STEP
	p.poolList[idx].Put(buffRef)
}
