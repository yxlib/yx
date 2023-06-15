// Copyright 2022 Guan Jianchang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yx

import (
	"errors"
	"sync/atomic"
)

var (
	ErrTryLockFail = errors.New("try lock failed")
)

type FastLock struct {
	state int32
}

func NewFastLock() *FastLock {
	return &FastLock{
		state: 0,
	}
}

func (l *FastLock) TryLock(maxCnt uint32) error {
	tryCnt := uint32(0)
	for {
		if atomic.CompareAndSwapInt32(&l.state, 0, 1) {
			// acquire OK
			return nil
		}

		if maxCnt == 0 {
			continue
		}

		tryCnt++
		if tryCnt >= maxCnt {
			break
		}
	}

	return ErrTryLockFail
}

func (l *FastLock) Unlock() {
	if ok := atomic.CompareAndSwapInt32(&l.state, 1, 0); !ok {
		panic("Unlock() failed")
	}
}
