// Copyright 2022 Guan Jianchang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yx

import (
	"errors"
	"sync"
)

var (
	ErrQueEmptyQueue = errors.New("empty queue")
)

//========================
//         Queue
//========================
type Queue interface {
	// Enqueue.
	// @param item, item to enqueue.
	// @return error, error.
	Enqueue(item interface{}) error

	// Dequeue.
	// @return interface{}, dequeue item.
	// @return error, error.
	Dequeue() (interface{}, error)

	// Get queue size.
	// @return uint64, queue size.
	GetSize() uint64
}

//========================
//    LinkedQueueNode
//========================
type LinkedQueueNode struct {
	Data interface{}
	Next *LinkedQueueNode
}

func NewLinkedQueueNode(data interface{}) *LinkedQueueNode {
	return &LinkedQueueNode{
		Data: data,
		Next: nil,
	}
}

//========================
//       LinkedQueue
//========================
type LinkedQueue struct {
	head *LinkedQueueNode
	back *LinkedQueueNode
	size uint64
}

func NewLinkedQueue() *LinkedQueue {
	return &LinkedQueue{
		head: nil,
		back: nil,
		size: 0,
	}
}

func (q *LinkedQueue) Enqueue(item interface{}) error {
	// if item == nil {
	// 	return errors.New("push nil item")
	// }

	n := NewLinkedQueueNode(item)

	if q.back == nil {
		q.back = n
		q.head = q.back
	} else {
		q.back.Next = n
		q.back = n
	}

	q.size++
	return nil
}

func (q *LinkedQueue) Dequeue() (interface{}, error) {
	if q.size == 0 {
		return nil, ErrQueEmptyQueue
	}

	n := q.head

	if q.head == q.back {
		q.head = nil
		q.back = nil
	} else {
		q.head = q.head.Next
	}

	n.Next = nil
	q.size--
	return n.Data, nil
}

func (q *LinkedQueue) GetSize() uint64 {
	return q.size
}

//========================
//    SyncLinkedQueue
//========================
type SyncLinkedQueue struct {
	lck *sync.Mutex
	que *LinkedQueue
}

func NewSyncLinkedQueue() *SyncLinkedQueue {
	return &SyncLinkedQueue{
		lck: &sync.Mutex{},
		que: NewLinkedQueue(),
	}
}

func (q *SyncLinkedQueue) Enqueue(item interface{}) error {
	q.lck.Lock()
	defer q.lck.Unlock()

	return q.que.Enqueue(item)
}

func (q *SyncLinkedQueue) Dequeue() (interface{}, error) {
	q.lck.Lock()
	defer q.lck.Unlock()

	return q.que.Dequeue()
}

func (q *SyncLinkedQueue) GetSize() uint64 {
	q.lck.Lock()
	defer q.lck.Unlock()

	return q.que.GetSize()
}
