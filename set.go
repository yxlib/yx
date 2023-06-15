// Copyright 2022 Guan Jianchang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yx

import "errors"

var (
	ErrSetKeyIsNil    = errors.New("key is nil")
	ErrSetNotIntKey   = errors.New("not int key")
	ErrSetNotUintKey  = errors.New("not uint key")
	ErrSetNotFloatKey = errors.New("not float key")
)

type SetType = uint8

const (
	SET_TYPE_INT SetType = iota
	SET_TYPE_UINT
	SET_TYPE_FLOAT
	SET_TYPE_OBJ
)

type Set interface {
	Add(data interface{}) error
	Remove(data interface{}) error
	Exist(data interface{}) (bool, error)
	GetElements() []interface{}
	GetSize() int
	Pop() (interface{}, bool)
}

func NewSet(t SetType) Set {
	if t == SET_TYPE_INT {
		return NewIntSet()
	} else if t == SET_TYPE_UINT {
		return NewUintSet()
	} else if t == SET_TYPE_FLOAT {
		return NewFloatSet()
	} else {
		return NewObjectSet()
	}
}

type IntSet struct {
	mapKv map[int64]bool
}

func NewIntSet() *IntSet {
	return &IntSet{
		mapKv: make(map[int64]bool),
	}
}

func (s *IntSet) Add(data interface{}) error {
	key, err := getIntKey(data)
	if err != nil {
		return err
	}

	_, ok := s.mapKv[key]
	if !ok {
		s.mapKv[key] = true
	}

	return nil
}

func (s *IntSet) Remove(data interface{}) error {
	key, err := getIntKey(data)
	if err != nil {
		return err
	}

	_, ok := s.mapKv[key]
	if ok {
		delete(s.mapKv, key)
	}

	return nil
}

func (s *IntSet) Exist(data interface{}) (bool, error) {
	key, err := getIntKey(data)
	if err != nil {
		return false, err
	}

	_, ok := s.mapKv[key]
	return ok, nil
}

func (s *IntSet) GetElements() []interface{} {
	objs := make([]interface{}, 0, len(s.mapKv))
	for k := range s.mapKv {
		objs = append(objs, k)
	}

	return objs
}

func (s *IntSet) GetSize() int {
	return len(s.mapKv)
}

func (s *IntSet) Pop() (interface{}, bool) {
	for k := range s.mapKv {
		delete(s.mapKv, k)
		return k, true
	}

	return nil, false
}

type UintSet struct {
	mapKv map[uint64]bool
}

func NewUintSet() *UintSet {
	return &UintSet{
		mapKv: make(map[uint64]bool),
	}
}

func (s *UintSet) Add(data interface{}) error {
	key, err := getUintKey(data)
	if err != nil {
		return err
	}

	_, ok := s.mapKv[key]
	if !ok {
		s.mapKv[key] = true
	}

	return nil
}

func (s *UintSet) Remove(data interface{}) error {
	key, err := getUintKey(data)
	if err != nil {
		return err
	}

	_, ok := s.mapKv[key]
	if ok {
		delete(s.mapKv, key)
	}

	return nil
}

func (s *UintSet) Exist(data interface{}) (bool, error) {
	key, err := getUintKey(data)
	if err != nil {
		return false, err
	}

	_, ok := s.mapKv[key]
	return ok, nil
}

func (s *UintSet) GetElements() []interface{} {
	objs := make([]interface{}, 0, len(s.mapKv))
	for k := range s.mapKv {
		objs = append(objs, k)
	}

	return objs
}

func (s *UintSet) GetSize() int {
	return len(s.mapKv)
}

func (s *UintSet) Pop() (interface{}, bool) {
	for k := range s.mapKv {
		delete(s.mapKv, k)
		return k, true
	}

	return nil, false
}

type FloatSet struct {
	mapKv map[float64]bool
}

func NewFloatSet() *FloatSet {
	return &FloatSet{
		mapKv: make(map[float64]bool),
	}
}

func (s *FloatSet) Add(data interface{}) error {
	key, err := getFloatKey(data)
	if err != nil {
		return err
	}

	_, ok := s.mapKv[key]
	if !ok {
		s.mapKv[key] = true
	}

	return nil
}

func (s *FloatSet) Remove(data interface{}) error {
	key, err := getFloatKey(data)
	if err != nil {
		return err
	}

	_, ok := s.mapKv[key]
	if ok {
		delete(s.mapKv, key)
	}

	return nil
}

func (s *FloatSet) Exist(data interface{}) (bool, error) {
	key, err := getFloatKey(data)
	if err != nil {
		return false, err
	}

	_, ok := s.mapKv[key]
	return ok, nil
}

func (s *FloatSet) GetElements() []interface{} {
	objs := make([]interface{}, 0, len(s.mapKv))
	for k := range s.mapKv {
		objs = append(objs, k)
	}

	return objs
}

func (s *FloatSet) GetSize() int {
	return len(s.mapKv)
}

func (s *FloatSet) Pop() (interface{}, bool) {
	for k := range s.mapKv {
		delete(s.mapKv, k)
		return k, true
	}

	return nil, false
}

type ObjectSet struct {
	mapKv map[interface{}]bool
}

func NewObjectSet() *ObjectSet {
	return &ObjectSet{
		mapKv: make(map[interface{}]bool),
	}
}

func (s *ObjectSet) Add(data interface{}) error {
	if data == nil {
		return ErrSetKeyIsNil
	}

	_, ok := s.mapKv[data]
	if !ok {
		s.mapKv[data] = true
	}

	return nil
}

func (s *ObjectSet) Remove(data interface{}) error {
	if data == nil {
		return ErrSetKeyIsNil
	}

	_, ok := s.mapKv[data]
	if ok {
		delete(s.mapKv, data)
	}

	return nil
}

func (s *ObjectSet) Exist(data interface{}) (bool, error) {
	if data == nil {
		return false, ErrSetKeyIsNil
	}

	_, ok := s.mapKv[data]
	return ok, nil
}

func (s *ObjectSet) GetElements() []interface{} {
	objs := make([]interface{}, 0, len(s.mapKv))
	for k := range s.mapKv {
		objs = append(objs, k)
	}

	return objs
}

func (s *ObjectSet) GetSize() int {
	return len(s.mapKv)
}

func (s *ObjectSet) Pop() (interface{}, bool) {
	for k := range s.mapKv {
		delete(s.mapKv, k)
		return k, true
	}

	return nil, false
}

func getIntKey(data interface{}) (int64, error) {
	var err error = nil
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
		err = ErrSetNotIntKey
	}

	return key, err
}

func getUintKey(data interface{}) (uint64, error) {
	var err error = nil
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
		err = ErrSetNotUintKey
	}

	return key, err
}

func getFloatKey(data interface{}) (float64, error) {
	var err error = nil
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
		err = ErrSetNotFloatKey
	}

	return key, err
}
