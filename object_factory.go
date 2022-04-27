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

type objectWorkshop struct {
	objType reflect.Type
	pool    *ObjectPool
}

func newObjectWorkshop(t reflect.Type, initReuseCnt uint64, maxCnt uint64) *objectWorkshop {
	w := &objectWorkshop{
		objType: t,
		pool:    NewObjectPool(maxCnt),
	}

	for i := 0; i < int(initReuseCnt); i++ {
		v := reflect.New(w.objType)
		obj := v.Interface()
		r, ok := obj.(Reuseable)
		if !ok {
			break
		}

		w.pushToPool(r)
	}

	return w
}

func (w *objectWorkshop) createObject() interface{} {
	obj, ok := w.popFromPool()
	if !ok {
		v := reflect.New(w.objType)
		obj = v.Interface()
	}

	return obj
}

func (w *objectWorkshop) reuseObject(v Reuseable) error {
	if v == nil {
		return ErrObjFactObjIsNil
	}

	t := reflect.TypeOf(v)
	t = t.Elem()
	if w.objType != t {
		return ErrObjFactNotTheSameType
	}

	return w.pushToPool(v)
}

func (w *objectWorkshop) popFromPool() (interface{}, bool) {
	return w.pool.Get()
}

func (w *objectWorkshop) pushToPool(v Reuseable) error {
	return w.pool.Reuse(v)
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
func (f *ObjectFactory) RegisterObject(obj Reuseable, initReuseCnt uint64, maxPoolCapacity uint64) (string, error) {
	name, err := GetClassReflectName(obj)
	if err != nil {
		return "", ErrObjFactObjIsNil
	}

	t := reflect.TypeOf(obj)
	err = f.createWorkshop(name, t, initReuseCnt, maxPoolCapacity)
	if err != nil {
		return "", err
	}

	return name, nil
}

// Get reflect type by the name.
// @param name, the name.
// @return reflect.Type, the reflect type.
// @return bool, true mean success, false mean failed.
func (f *ObjectFactory) GetReflectType(name string) (reflect.Type, bool) {
	f.lckWs.RLock()
	defer f.lckWs.RUnlock()

	w, ok := f.mapName2Workshop[name]
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
	f.lckWs.RLock()
	defer f.lckWs.RUnlock()

	w, ok := f.mapName2Workshop[name]
	if !ok {
		return nil, ErrObjFactWsNotExist
	}

	return w.createObject(), nil
}

// reuse an object by the name.
// @param v, the object.
// @param name, the name.
// @return error, the error.
func (f *ObjectFactory) ReuseObject(v Reuseable, name string) error {
	f.lckWs.RLock()
	defer f.lckWs.RUnlock()

	w, ok := f.mapName2Workshop[name]
	if !ok {
		return ErrObjFactWsNotExist
	}

	return w.reuseObject(v)
}

func (f *ObjectFactory) createWorkshop(name string, t reflect.Type, initReuseCnt uint64, maxPoolCapacity uint64) error {
	f.lckWs.Lock()
	defer f.lckWs.Unlock()

	_, ok := f.mapName2Workshop[name]
	if ok {
		return ErrObjFactWsExist
	}

	f.mapName2Workshop[name] = newObjectWorkshop(t, initReuseCnt, maxPoolCapacity)
	return nil
}
