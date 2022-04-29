// Copyright 2022 Guan Jianchang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yx

import (
	"errors"
	"reflect"
	"sync"
)

var (
	ErrObjFactObjIsNil       = errors.New("object is nil")
	ErrObjFactWsExist        = errors.New("object workshop is exist")
	ErrObjFactWsNotExist     = errors.New("object workshop is not exist")
	ErrObjFactNotTheSameType = errors.New("not the same type")
)

// type Reuseable interface {
// 	Reset()
// }

type objectWorkshop struct {
	objType reflect.Type
	pool    *sync.Pool
	queue   Queue
}

func newObjectWorkshop(t reflect.Type, newFunc func() interface{}, maxCnt uint64) *objectWorkshop {
	w := &objectWorkshop{
		objType: t,
		pool:    nil,
	}

	if newFunc != nil {
		w.pool = &sync.Pool{
			New: newFunc,
		}
	} else {
		w.queue = NewSyncLimitQueue(maxCnt)
	}

	return w
}

func (w *objectWorkshop) createObject() interface{} {
	if w.pool != nil {
		return w.pool.Get()
	}

	obj, err := w.queue.Dequeue()
	if err != nil {
		v := reflect.New(w.objType)
		obj = v.Interface()
	}

	return obj
}

func (w *objectWorkshop) reuseObject(v interface{}) error {
	if v == nil {
		return ErrObjFactObjIsNil
	}

	if w.pool != nil {
		w.pool.Put(v)
		return nil
	}

	return w.queue.Enqueue(v)
}

//=========================
//      ObjectFactory
//=========================
type ObjectFactory struct {
	mapName2Workshop map[string]*objectWorkshop
	lckWs            *sync.RWMutex
}

func NewObjectFactory() *ObjectFactory {
	return &ObjectFactory{
		mapName2Workshop: make(map[string]*objectWorkshop),
		lckWs:            &sync.RWMutex{},
	}
}

// Register object.
// @param obj, the reuseabe object.
// @param initReuseCnt, the init count in the pool.
// @param maxPoolCapacity, the max count in the pool.
// @return error, the error.
func (f *ObjectFactory) RegisterObject(obj interface{}, newFunc func() interface{}, maxPoolCapacity uint64) (string, error) {
	if obj == nil {
		return "", ErrObjFactObjIsNil
	}

	t := reflect.TypeOf(obj)
	t = t.Elem()
	name := GetReflectNameByType(t)
	f.createWorkshop(name, t, newFunc, maxPoolCapacity)
	return name, nil
}

// Get reflect type by the name.
// @param name, the name.
// @return reflect.Type, the reflect type.
// @return bool, true mean success, false mean failed.
func (f *ObjectFactory) GetReflectType(name string) (reflect.Type, bool) {
	w, ok := f.getWorkshop(name)
	if ok {
		return w.objType, ok
	}

	return nil, false
}

// create an object by the name.
// @param name, the name.
// @return interface{}, the object is created.
// @return error, the error.
func (f *ObjectFactory) CreateObject(name string) (interface{}, error) {
	w, ok := f.getWorkshop(name)
	if !ok {
		return nil, ErrObjFactWsNotExist
	}

	return w.createObject(), nil
}

// reuse an object by the name.
// @param v, the object.
// @param name, the name.
// @return error, the error.
func (f *ObjectFactory) ReuseObject(v interface{}, name string) error {
	w, ok := f.getWorkshop(name)
	if !ok {
		return ErrObjFactWsNotExist
	}

	return w.reuseObject(v)
}

func (f *ObjectFactory) createWorkshop(name string, t reflect.Type, newFunc func() interface{}, maxPoolCapacity uint64) {
	f.lckWs.Lock()
	defer f.lckWs.Unlock()

	_, ok := f.mapName2Workshop[name]
	if !ok {
		f.mapName2Workshop[name] = newObjectWorkshop(t, newFunc, maxPoolCapacity)
	}
}

func (f *ObjectFactory) getWorkshop(name string) (*objectWorkshop, bool) {
	f.lckWs.RLock()
	defer f.lckWs.RUnlock()

	w, ok := f.mapName2Workshop[name]
	return w, ok
}
