// Copyright 2022 Guan Jianchang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yx

type SetType = uint8

const (
	SET_TYPE_INT SetType = iota
	SET_TYPE_UINT
	SET_TYPE_FLOAT
	SET_TYPE_OBJ
)

type Set struct {
	t        SetType
	intSet   map[int64]bool
	uintSet  map[uint64]bool
	floatSet map[float64]bool
	objSet   map[interface{}]bool
}

func NewSet(t SetType) *Set {
	s := &Set{
		t: t,
	}

	if t == SET_TYPE_INT {
		s.intSet = make(map[int64]bool)
	} else if t == SET_TYPE_UINT {
		s.uintSet = make(map[uint64]bool)
	} else if t == SET_TYPE_FLOAT {
		s.floatSet = make(map[float64]bool)
	} else {
		s.objSet = make(map[interface{}]bool)
	}

	return s
}

// Add item.
// @param data, the item to add
func (s *Set) Add(data interface{}) {
	switch s.t {
	case SET_TYPE_INT:
		s.addIntVal(data)

	case SET_TYPE_UINT:
		s.addUintVal(data)

	case SET_TYPE_FLOAT:
		s.addFloatVal(data)

	default:
		s.addObjVal(data)
	}
}

// Remove item.
// @param data, the item to remove
func (s *Set) Remove(data interface{}) {
	switch s.t {
	case SET_TYPE_INT:
		s.removeIntVal(data)

	case SET_TYPE_UINT:
		s.removeUintVal(data)

	case SET_TYPE_FLOAT:
		s.removeFloatVal(data)

	default:
		s.removeObjVal(data)
	}
}

// Is the item exist.
// @return bool, true exist, false not exist.
func (s *Set) Exist(data interface{}) bool {
	switch s.t {
	case SET_TYPE_INT:
		return s.existIntVal(data)

	case SET_TYPE_UINT:
		return s.existUintVal(data)

	case SET_TYPE_FLOAT:
		return s.existFloatVal(data)

	default:
		return s.existObjVal(data)
	}
}

// Get all items.
// @return []interface{}, all item array.
func (s *Set) GetElements() []interface{} {
	switch s.t {
	case SET_TYPE_INT:
		return s.getIntElements()

	case SET_TYPE_UINT:
		return s.getUintElements()

	case SET_TYPE_FLOAT:
		return s.getFloatElements()

	default:
		return s.getObjElements()
	}
}

// Get set size.
// @return int, the set size.
func (s *Set) GetSize() int {
	switch s.t {
	case SET_TYPE_INT:
		return len(s.intSet)

	case SET_TYPE_UINT:
		return len(s.uintSet)

	case SET_TYPE_FLOAT:
		return len(s.floatSet)

	default:
		return len(s.objSet)
	}
}

// Pop an items.
// @return interface{}, a random item.
// @return bool, true mean success, false mean failed.
func (s *Set) Pop() (interface{}, bool) {
	switch s.t {
	case SET_TYPE_INT:
		return s.popInt()

	case SET_TYPE_UINT:
		return s.popUint()

	case SET_TYPE_FLOAT:
		return s.popFloat()

	default:
		return s.popObj()
	}
}

func (s *Set) addIntVal(data interface{}) {
	key := s.getIntKey(data)
	_, ok := s.intSet[key]
	if !ok {
		s.intSet[key] = true
	}
}

func (s *Set) removeIntVal(data interface{}) {
	key := s.getIntKey(data)
	_, ok := s.intSet[key]
	if ok {
		delete(s.intSet, key)
	}
}

func (s *Set) existIntVal(data interface{}) bool {
	key := s.getIntKey(data)
	_, ok := s.intSet[key]
	return ok
}

func (s *Set) getIntElements() []interface{} {
	objs := make([]interface{}, 0, len(s.intSet))
	for k := range s.intSet {
		objs = append(objs, k)
	}

	return objs
}

func (s *Set) popInt() (interface{}, bool) {
	for k := range s.intSet {
		delete(s.intSet, k)
		return k, true
	}

	return nil, false
}

func (s *Set) addUintVal(data interface{}) {
	key := s.getUintKey(data)
	_, ok := s.uintSet[key]
	if !ok {
		s.uintSet[key] = true
	}
}

func (s *Set) removeUintVal(data interface{}) {
	key := s.getUintKey(data)
	_, ok := s.uintSet[key]
	if ok {
		delete(s.uintSet, key)
	}
}

func (s *Set) existUintVal(data interface{}) bool {
	key := s.getUintKey(data)
	_, ok := s.uintSet[key]
	return ok
}

func (s *Set) getUintElements() []interface{} {
	objs := make([]interface{}, 0, len(s.uintSet))
	for k := range s.uintSet {
		objs = append(objs, k)
	}

	return objs
}

func (s *Set) popUint() (interface{}, bool) {
	for k := range s.uintSet {
		delete(s.uintSet, k)
		return k, true
	}

	return nil, false
}

func (s *Set) addFloatVal(data interface{}) {
	key := s.getFloatKey(data)
	_, ok := s.floatSet[key]
	if !ok {
		s.floatSet[key] = true
	}
}

func (s *Set) removeFloatVal(data interface{}) {
	key := s.getFloatKey(data)
	_, ok := s.floatSet[key]
	if ok {
		delete(s.floatSet, key)
	}
}

func (s *Set) existFloatVal(data interface{}) bool {
	key := s.getFloatKey(data)
	_, ok := s.floatSet[key]
	return ok
}

func (s *Set) getFloatElements() []interface{} {
	objs := make([]interface{}, 0, len(s.floatSet))
	for k := range s.floatSet {
		objs = append(objs, k)
	}

	return objs
}

func (s *Set) popFloat() (interface{}, bool) {
	for k := range s.floatSet {
		delete(s.floatSet, k)
		return k, true
	}

	return nil, false
}

func (s *Set) addObjVal(data interface{}) {
	if data == nil {
		return
	}

	_, ok := s.objSet[data]
	if !ok {
		s.objSet[data] = true
	}
}

func (s *Set) removeObjVal(data interface{}) {
	if data == nil {
		return
	}

	_, ok := s.objSet[data]
	if ok {
		delete(s.objSet, data)
	}
}

func (s *Set) existObjVal(data interface{}) bool {
	if data == nil {
		return false
	}

	_, ok := s.objSet[data]
	return ok
}

func (s *Set) getObjElements() []interface{} {
	objs := make([]interface{}, 0, len(s.objSet))
	for k := range s.objSet {
		objs = append(objs, k)
	}

	return objs
}

func (s *Set) popObj() (interface{}, bool) {
	for k := range s.objSet {
		delete(s.objSet, k)
		return k, true
	}

	return nil, false
}

func (s *Set) getIntKey(data interface{}) int64 {
	key := int64(0)

	switch val := data.(type) {
	case int8:
		key = int64(val)
	case int16:
		key = int64(val)
	case int:
		key = int64(val)
	case int32:
		key = int64(val)
	case int64:
		key = val
	case uint8:
		key = int64(val)
	case uint16:
		key = int64(val)
	case uint:
		key = int64(val)
	case uint32:
		key = int64(val)
	case uint64:
		key = int64(val)
	case float32:
		key = int64(val)
	case float64:
		key = int64(val)
	default:
	}

	return key
}

func (s *Set) getUintKey(data interface{}) uint64 {
	key := uint64(0)

	switch val := data.(type) {
	case int8:
		key = uint64(val)
	case int16:
		key = uint64(val)
	case int:
		key = uint64(val)
	case int32:
		key = uint64(val)
	case int64:
		key = uint64(val)
	case uint8:
		key = uint64(val)
	case uint16:
		key = uint64(val)
	case uint:
		key = uint64(val)
	case uint32:
		key = uint64(val)
	case uint64:
		key = val
	case float32:
		key = uint64(val)
	case float64:
		key = uint64(val)
	default:
	}

	return key
}

func (s *Set) getFloatKey(data interface{}) float64 {
	key := float64(0)

	switch val := data.(type) {
	case int8:
		key = float64(val)
	case int16:
		key = float64(val)
	case int:
		key = float64(val)
	case int32:
		key = float64(val)
	case int64:
		key = float64(val)
	case uint8:
		key = float64(val)
	case uint16:
		key = float64(val)
	case uint:
		key = float64(val)
	case uint32:
		key = float64(val)
	case uint64:
		key = float64(val)
	case float32:
		key = float64(val)
	case float64:
		key = val
	default:
	}

	return key
}
