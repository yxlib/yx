// Copyright 2022 Guan Jianchang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yx

import "errors"

var (
	ErrObjPoolObjNil = errors.New("object is nil")
)

type Reuseable interface {
	Reset()
}

type ObjectPool struct {
	queue  Queue
	maxCnt uint64
}

func NewObjectPool(maxCnt uint64) *ObjectPool {
	return &ObjectPool{
		queue:  NewSyncLimitQueue(maxCnt),
		maxCnt: maxCnt,
	}
}

// Get an object.
// @return Reuseable, the object.
// @return bool, true mean success, false mean failed.
func (p *ObjectPool) Get() (Reuseable, bool) {
	item, err := p.queue.Dequeue()
	if err != nil {
		return nil, false
	}

	obj, ok := item.(Reuseable)
	if !ok {
		return nil, false
	}

	return obj, true
}

// Reuse an object.
// @param obj, the object.
// @return error, the error.
func (p *ObjectPool) Reuse(obj Reuseable) error {
	if obj == nil {
		return ErrObjPoolObjNil
	}

	obj.Reset()
	err := p.queue.Enqueue(obj)
	return err
}
