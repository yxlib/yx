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
	C         chan byte
	chanClose chan byte
}

func NewEvent() *Event {
	return &Event{
		C:         make(chan byte, 1),
		chanClose: make(chan byte),
	}
}

// Send event.
func (e *Event) Send() error {
	select {
	case <-e.chanClose:
		return ErrEvtClosed

	default:
		select {
		case e.C <- 1:
		default:
		}
	}

	return nil
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
	select {
	case <-e.chanClose:
		return ErrEvtClosed
	case <-e.C:
	}

	return nil
	// _, ok := <-e.C
	// if ok {
	// 	return nil
	// }

	// return ErrEvtClosed
}

// Close the event.
func (e *Event) Close() {
	select {
	case <-e.chanClose:
		return
	default:
		close(e.chanClose)
	}
	// if !e.bClosed {
	// 	e.bClosed = true
	// 	close(e.C)
	// }
}

// Wait event until timeout.
// @param timeoutMSec, timeout after millisecond. 0 mean check at once, only can use in one waiter scene.
// @return error, ErrEvtWaitTimeout mean timeout.
func (e *Event) WaitUntilTimeout(timeoutMSec uint32) error {
	if timeoutMSec == 0 {
		select {
		case <-e.C:
			return nil
		default:
			return ErrEvtWaitTimeout
		}
	}

	var err error = nil
	t := time.NewTimer(time.Millisecond * time.Duration(timeoutMSec))

	select {
	case <-e.chanClose:
		t.Stop()
		err = ErrEvtClosed

	case <-e.C:
		t.Stop()

	case <-t.C:
		err = ErrEvtWaitTimeout
	}

	return err
}
