// Copyright 2022 Guan Jianchang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yx

import (
	"errors"
	"time"
)

var (
	ErrEvtWaitTimeout = errors.New("event wait time out")
	ErrEvtClosed      = errors.New("event closed")
)

type Event struct {
	chanBroadcast chan byte
	lck           *FastLock
	// lck           *sync.Mutex
	// chanClose chan byte
}

func NewEvent() *Event {
	return &Event{
		chanBroadcast: make(chan byte, 1),
		lck:           NewFastLock(),
		// lck:           &sync.Mutex{},
		// C:         make(chan byte, 1),
		// chanClose: make(chan byte),
	}
}

// Broadcast event.
func (e *Event) Broadcast() error {
	// e.lck.Lock()
	if err := e.lck.TryLock(0); err != nil {
		return err
	}

	ch := e.chanBroadcast
	if ch != nil {
		e.chanBroadcast = make(chan byte, 1)
	}

	e.lck.Unlock()

	if ch != nil {
		close(ch)
	}
	return nil

	// select {
	// case <-e.chanClose:
	// 	return ErrEvtClosed

	// default:
	// 	select {
	// 	case e.C <- 1:
	// 	default:
	// 	}
	// }

	// return nil
	// if e.bClosed {
	// 	return ErrEvtClosed
	// }

	// if len(e.C) == 0 {
	// 	e.C <- 1
	// }

	// return nil
}

// Wait event.
func (e *Event) Wait() error {
	// e.lck.Lock()
	if err := e.lck.TryLock(0); err != nil {
		return err
	}

	ch := e.chanBroadcast
	e.lck.Unlock()

	if ch == nil {
		return nil
	}

	<-ch
	return nil

	// select {
	// case <-e.chanClose:
	// 	return ErrEvtClosed
	// case <-e.C:
	// }

	// return nil
	// _, ok := <-e.C
	// if ok {
	// 	return nil
	// }

	// return ErrEvtClosed
}

func (e *Event) GetChan() chan byte {
	// e.lck.Lock()
	if e.lck.TryLock(0) != nil {
		return nil
	}

	ch := e.chanBroadcast
	e.lck.Unlock()

	return ch
}

// Close the event.
func (e *Event) Close() {
	// e.lck.Lock()
	if e.lck.TryLock(0) != nil {
		return
	}

	ch := e.chanBroadcast
	e.chanBroadcast = nil
	e.lck.Unlock()

	if ch != nil {
		close(ch)
	}

	// select {
	// case <-e.chanClose:
	// 	return
	// default:
	// 	close(e.chanClose)
	// }
	// if !e.bClosed {
	// 	e.bClosed = true
	// 	close(e.C)
	// }
}

func (e *Event) IsClose() bool {
	return e.GetChan() == nil
}

// Wait event until timeout.
// @param timeoutMSec, timeout after millisecond. 0 mean check at once, only can use in one waiter scene.
// @return error, ErrEvtWaitTimeout mean timeout.
func (e *Event) WaitUntilTimeout(timeoutMSec uint32) error {
	ch := e.GetChan()
	if ch == nil {
		return ErrEvtClosed
	}

	if timeoutMSec == 0 {
		select {
		case <-ch:
			return nil
		default:
			return ErrEvtWaitTimeout
		}
	}

	var err error = nil
	t := time.NewTimer(time.Millisecond * time.Duration(timeoutMSec))

	select {
	case <-ch:
		t.Stop()

	case <-t.C:
		err = ErrEvtWaitTimeout
	}

	return err
}
