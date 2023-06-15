// Copyright 2022 Guan Jianchang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yx

import (
	"errors"
)

var (
	ErrIdGenIdUseOut = errors.New("id use out")
)

type IdGenerator struct {
	// lck      *sync.Mutex
	lck      *FastLock
	curId    uint64
	maxId    uint64
	reuseIds *UintSet
}

func NewIdGenerator(min uint64, max uint64) *IdGenerator {
	return &IdGenerator{
		// lck:      &sync.Mutex{},
		lck:      NewFastLock(),
		curId:    min,
		maxId:    max,
		reuseIds: NewUintSet(),
	}
}

// Assign an id.
// @return uint64, the assign id.
// @return error, ErrIdGenIdUseOut mean use out.
func (g *IdGenerator) GetId() (uint64, error) {
	// g.lck.Lock()
	if err := g.lck.TryLock(0); err != nil {
		return 0, err
	}

	defer g.lck.Unlock()

	id, ok := g.assignId()
	if ok {
		return id, nil
	}

	id, err := g.popReuseId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Reuse an id.
// @param id, the reuse id.
func (g *IdGenerator) ReuseId(id uint64) {
	if g.lck.TryLock(0) != nil {
		return
	}

	defer g.lck.Unlock()

	g.pushReuseId(id)
}

func (g *IdGenerator) pushReuseId(id uint64) {
	g.reuseIds.Add(id)
}

func (g *IdGenerator) popReuseId() (uint64, error) {
	id, ok := g.reuseIds.Pop()
	if !ok {
		return 0, ErrIdGenIdUseOut
	}

	return id.(uint64), nil
}

func (g *IdGenerator) assignId() (uint64, bool) {
	if g.curId <= g.maxId {
		id := g.curId
		g.curId++
		return id, true
	}

	return 0, false
}
