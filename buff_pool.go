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
	poolList  []*BuffPool
	minSize   uint32
	maxSize   uint32
	queueStep uint32
}

func NewBuffFactory(minSize uint32, maxSize uint32, queueStep uint32) *BuffFactory {
	if minSize == 0 {
		minSize = BP_MIN_BUFF_SIZE
	}

	if maxSize == 0 {
		maxSize = BP_MAX_BUFF_SIZE
	}

	if queueStep == 0 {
		queueStep = BP_QUEUE_STEP
	}

	if minSize > maxSize {
		minSize, maxSize = maxSize, minSize
	}

	listCnt := (maxSize-minSize)/queueStep + 1
	p := &BuffFactory{
		poolList:  make([]*BuffPool, listCnt),
		minSize:   minSize,
		maxSize:   maxSize,
		queueStep: queueStep,
	}

	for i := uint32(0); i < listCnt; i++ {
		size := minSize + i*queueStep
		p.poolList[i] = NewBuffPool(size)
	}

	return p
}

func (p *BuffFactory) CreateBuff(size uint32) *[]byte {
	if size > uint32(p.maxSize) {
		buff := make([]byte, size)
		return &buff
	}

	idx := uint32(0)
	if size > p.minSize {
		idx = (size - p.minSize) / p.queueStep
		if (size-p.minSize)%p.queueStep != 0 {
			idx++
		}
	}

	buffRef := p.poolList[idx].Get()
	return buffRef
}

func (p *BuffFactory) ReuseBuff(buffRef *[]byte) {
	size := uint32(len(*buffRef))
	if size < p.minSize || size > p.maxSize {
		return
	}

	if (size-p.minSize)%p.queueStep != 0 {
		return
	}

	idx := (size - p.minSize) / p.queueStep
	p.poolList[idx].Put(buffRef)
}
