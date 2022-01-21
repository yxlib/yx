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
	C       chan byte
	bClosed bool
}

func NewEvent() *Event {
	return &Event{
		C:       make(chan byte, 1),
		bClosed: false,
	}
}

// Send event.
func (e *Event) Send() error {
	if e.bClosed {
		return ErrEvtClosed
	}

	if len(e.C) == 0 {
		e.C <- 1
	}

	return nil
}

// Wait event.
func (e *Event) Wait() error {
	_, ok := <-e.C
	if ok {
		return nil
	}

	return ErrEvtClosed
}

// Close the event.
func (e *Event) Close() {
	if !e.bClosed {
		e.bClosed = true
		close(e.C)
	}
}

// Wait event until timeout.
// @param timeoutMSec, timeout after millisecond. 0 mean check at once, only can use in one waiter scene.
// @return error, ErrEvtWaitTimeout mean timeout.
func (e *Event) WaitUntilTimeout(timeoutMSec uint32) error {
	if e.bClosed {
		return ErrEvtClosed
	}

	// check at once
	if timeoutMSec == 0 {
		if len(e.C) == 1 {
			<-e.C
			return nil
		}

		return ErrEvtWaitTimeout
	}

	// wait until timeout
	var err error = nil
	t := time.NewTimer(time.Millisecond * time.Duration(timeoutMSec))

	select {
	case <-t.C:
		err = ErrEvtWaitTimeout

	case _, ok := <-e.C:
		t.Stop()
		if !ok {
			err = ErrEvtClosed
		}
	}

	return err
}
